package models

type Method string

const (
	PostTopic   Method = "post_topic"
	DeleteTopic Method = "delete_topic"
	ListTopics  Method = "list_topics"
	PostLink    Method = "post_link"
	PickLink    Method = "pick_link"
	ListLinks   Method = "list_links"
	DeleteLink  Method = "delete_links"
	Help        Method = "help"
)

type Arg string

const (
	Alias Arg = "alias"
	Link  Arg = "link"
	Topic Arg = "topic"
)
