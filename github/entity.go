/*
Copyright © 2024 jamie HERE <EMAIL ADDRESS>
*/
package github

type Repository struct {
	Id          int    `json:"id"`
	NodeId      string `json:"node_id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Private     bool   `json:"private"`
	HtmlUrl     string `json:"html_url"`
	Description string `json:"description"`
	Fork        bool   `json:"fork"`
	Url         string `json:"url"`
}
