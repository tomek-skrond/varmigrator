package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/jefflinse/githubsecret"
)

func NewVarDB(repo_data *RepoData) (*VarDB, error) {
	return &VarDB{
		repo_data: repo_data,
	}, nil
}

func (v *VarDB) GetVars() error {
	client := http.Client{}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/actions/variables",
			v.repo_data.Username,
			v.repo_data.Repo),
		nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// Set request headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+v.repo_data.GithubToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(body))

	var repoVars VarResponse
	if err := json.NewDecoder(resp.Body).Decode(&repoVars); err != nil {
		return err
	}

	for i, va := range repoVars.Variables {
		va.Id = i
	}

	v.vars = repoVars.Variables

	return nil
}

type PublicKeyResponse struct {
	KeyID     string `json:"key_id"`
	PublicKey string `json:"key"`
}

func (v *VarDB) GetRepoPublicKey() (string, string, error) {
	fmt.Println(v.repo_data.Username, v.repo_data.Repo)
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/secrets/public-key", v.repo_data.Username, v.repo_data.Repo)

	client, req, err := v.BuildGithubApiCall("GET", url)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	// var responseBody []byte
	// if resp.Body != nil {
	// 	responseBody, _ = io.ReadAll(resp.Body)
	// }
	// log.Printf("Error decoding response body: %s\nResponse body: %s", "", responseBody)

	fmt.Println(resp.StatusCode)
	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("error: %s", resp.Status)
	}

	// Decode the response body
	var publicKeyResp *PublicKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&publicKeyResp); err != nil {
		var responseBody []byte
		if resp.Body != nil {
			responseBody, _ = ioutil.ReadAll(resp.Body)
		}
		log.Printf("Error decoding response body: %s\nResponse body: %s", err, responseBody)
		return "", "", fmt.Errorf("error decoding response body: %s", err)
	}

	fmt.Println(publicKeyResp.PublicKey)
	// Return the public key
	return publicKeyResp.PublicKey, publicKeyResp.KeyID, nil
}

type apiConfigFunc func(req *http.Request) *http.Request

func (v *VarDB) configureSecretEncryption(req *http.Request) *http.Request {
	publicKey, keyID, err := v.GetRepoPublicKey()
	if err != nil {
		return nil
	}

	encryptedSecret, err := githubsecret.Encrypt(publicKey, v.currentSecret)
	if err != nil {
		fmt.Println("Encrypting error:", err)
		return req
	}

	requestBodyJSON, err := json.Marshal(&EncryptedRequestBody{
		EncryptedValue: encryptedSecret,
		KeyID:          keyID,
	})
	if err != nil {
		// Handle error
		fmt.Println("adding request body error:", err)
		return req
	}

	req.Body = io.NopCloser(bytes.NewBuffer([]byte(requestBodyJSON)))

	return req

}
func (v *VarDB) BuildGithubApiCall(method string, url string, config ...apiConfigFunc) (*http.Client, *http.Request, error) {

	client := http.Client{}
	req, err := http.NewRequest(
		method,
		url,
		nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return &client, nil, err
	}

	// Set request headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+v.repo_data.GithubToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	for _, c := range config {
		req = c(req)
	}

	return &client, req, nil
}
func (v *VarDB) GetSecrets() error {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/secrets", v.repo_data.Username, v.repo_data.Repo)

	client, req, err := v.BuildGithubApiCall("GET", url)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}

	var repoSecrets SecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&repoSecrets); err != nil {
		return err
	}

	for i, s := range repoSecrets.Secrets {
		s.Id = i
	}
	v.secrets = repoSecrets.Secrets

	return nil
}

func (v *VarDB) PrintAll() {
	fmt.Println("All available project secrets: ")
	for i, variable := range v.vars {
		fmt.Printf("%d. %s (current: %s)\n", i+1, variable.Name, variable.Value)
	}
	fmt.Println("")

	fmt.Println("All available project variables: ")
	for i, secret := range v.secrets {
		fmt.Printf("%d. %s\n", i+1, secret.Name)
	}
	fmt.Println("")
}

func (v *VarDB) DetermineVariableToEdit(userInput string) error {
	pattern := `^(v|s)(\d+)$`
	re := regexp.MustCompile(pattern)

	if re.MatchString(userInput) {

		matches := re.FindStringSubmatch(userInput)

		if len(matches) == 3 {
			varType := matches[1]
			varIndex, err := strconv.Atoi(matches[2])
			if err != nil {
				return err
			}

			if varType == "s" && len(v.secrets) >= varIndex {
				editContent := GetInput(EditSecretMessages)
				// fmt.Println(editContent)
				v.EditSecret(varIndex-1, editContent)
			}
			if varType == "v" && len(v.secrets) >= varIndex {
				editContent := GetInput(EditVariableMessages)
				v.EditVar(varIndex-1, editContent)
			}
		} else {
			fmt.Println("No matches found")
		}
		return nil
	} else {
		fmt.Println("your input does not match the pattern")
		return fmt.Errorf("your input does not match the pattern")
	}
}
func (v *VarDB) EditSecret(index int, editContent string) error {
	v.currentSecret = editContent
	secretName := v.secrets[index].Name

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/secrets/%s", v.repo_data.Username, v.repo_data.Repo, secretName)

	client, req, err := v.BuildGithubApiCall("PUT", url, v.configureSecretEncryption)
	if err != nil {
		return fmt.Errorf("error building API call: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("status code of changing secret: ", resp.StatusCode)
	if resp.StatusCode != 204 {
		return fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	return nil
}

func (v *VarDB) EditVar(index int, editContent string) error {
	varName := v.vars[index].Name
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/variables/%s", v.repo_data.Username, v.repo_data.Repo, varName)
	client, req, err := v.BuildGithubApiCall("PATCH", url)
	if err != nil {
		return fmt.Errorf("error building API call: %v", err)
	}

	requestBodyJSON, err := json.Marshal(map[string]string{
		"name":  varName,
		"value": editContent,
	})
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(requestBodyJSON))
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(body))

	fmt.Println("status code of changing var: ", resp.StatusCode)
	if resp.StatusCode != 204 {
		return fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}
	return nil
}
