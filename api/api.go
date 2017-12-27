package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Request struct {
	RequestMethod string
	UserAuth      bool

	Resource string
	Method   string

	MethodParams []Param
	ApiParams    []Param
	PostParams   []Param
}

var Scheme string = "https"

func (r Request) BuildPath() string {
	var parts []string
	for _, v := range r.ApiParams {
		parts = append(parts, v.Name, v.Value)
	}
	apiParams := strings.Join(parts, ",")

	parts = []string{r.Resource, r.Method}
	for _, v := range r.MethodParams {
		parts = append(parts, v.Value)
	}
	if len(apiParams) > 0 {
		parts = append(parts, apiParams)
	}
	return strings.Join(parts, "/")
}

func (r Request) BuildURL() url.URL {
	return url.URL{
		Scheme: Scheme,
		Host:   "a.wykop.pl",
		Path:   r.BuildPath(),
	}
}

var ErrWrongRequest = errors.New("wrong request")

func (r Request) Build() (*http.Request, error) {
	if r.Resource == "" || r.RequestMethod == "" || r.Method == "" {
		return nil, ErrWrongRequest
	}
	u := r.BuildURL()
	return http.NewRequest(r.RequestMethod, u.String(), nil)
}

type Response struct{}

type Param struct {
	Name  string
	Value string
}

type Client struct {
	appKey     string
	userKey    string
	httpClient *http.Client
}

// New creates a new client API. The transport is passed to http.Client.
func New(appKey, userKey string, transport http.RoundTripper) *Client {
	return &Client{appKey, userKey, &http.Client{Transport: transport}}
}

func (c *Client) Do(r *Request, resp interface{}) error {
	r.ApiParams = append(r.ApiParams, Param{"appkey", c.appKey})
	if r.UserAuth {
		r.ApiParams = append(r.ApiParams, Param{"userkey", c.userKey})
	}
	req, err := r.Build()
	if err != nil {
		return err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if resp != nil {
		//r := bufio.NewReader(res.Body)
		//rune, _, _ := r.ReadRune()
		//r.UnreadRune()
		//if rune == '{' {
		//
		//}
		dec := json.NewDecoder(res.Body)
		return dec.Decode(resp)
	}
	return nil
}

type PromotedSort string

const (
	PromotedByDay   (PromotedSort) = "day"
	PromotedByWeek  (PromotedSort) = "week"
	PromotedByMonth (PromotedSort) = "month"
)

type UpcomingSort string

const (
	UpcomingDate     (UpcomingSort) = "date"     // (najnowsze),
	UpcomingVotes    (UpcomingSort) = "votes"    // (wykopywane),
	UpcomingComments (UpcomingSort) = "comments" // comments (komentowane)
)

type LinksResource interface {
	Promoted(page int, sort PromotedSort) ([]Link, error)
	Upcoming(page int, sort UpcomingSort) ([]Link, error)
}

type resLinks struct {
	c *Client
}

func (r resLinks) Promoted(page int, sort PromotedSort) ([]Link, error) {
	var links []Link
	req := &Request{
		RequestMethod: http.MethodGet,
		Resource:      "links",
		Method:        "promoted",
		ApiParams:     []Param{{"page", strconv.Itoa(page)}, {"sort", string(sort)}},
	}
	err := r.c.Do(req, &links)
	return links, err
}

func (r resLinks) Upcoming(page int, sort UpcomingSort) ([]Link, error) {
	var links []Link
	req := &Request{
		RequestMethod: http.MethodGet,
		Resource:      "links",
		Method:        "upcoming",
		ApiParams:     []Param{{"page", strconv.Itoa(page)}, {"sort", string(sort)}},
	}
	err := r.c.Do(req, &links)
	return links, err
}

func (r resLinks) Index(id int) (*Link, error) {
	var link *Link
	req := &Request{
		UserAuth:      true,
		RequestMethod: http.MethodGet,
		Resource:      "link",
		Method:        "index",
		MethodParams:  []Param{{"param1", strconv.Itoa(id)}},
	}
	err := r.c.Do(req, &link)
	return link, err
}

func (r resLinks) Comments(id int) ([]*Comment, error) {
	panic("implement me")
}

func (r resLinks) Reports(id int) ([]*Bury, error) {
	panic("implement me")
}

func (r resLinks) Digs(id int) ([]*Dig, error) {
	var digs []*Dig
	req := &Request{
		RequestMethod: http.MethodGet,
		Resource:      "link",
		Method:        "digs",
		MethodParams:  []Param{{"param1", strconv.Itoa(id)}},
	}
	err := r.c.Do(req, &digs)
	return digs, err
}

func (r resLinks) Related(id int) ([]*RelatedLink, error) {
	panic("implement me")
}

func (r resLinks) BuryReasons() {
	panic("implement me")
}

func (c *Client) Links() LinksResource {
	return resLinks{c}
}

func (c *Client) Link() LinkResource {
	return resLinks{c}
}

/*
	LINK

	Pole	Wartość	Typ
	id	identyfikator linka	int
	title	tytuł linka	string
	description	opis	string
	tags	tagi	string
	url	adres url w serwisie wykop.pl	uri
	source_url	adres źródłowy	uri
	vote_count	liczba głosów	int
	comment_count	liczba komentarzy	int
	report_count	liczba zakopów	int
	date	data dodania	date
	author	login dodającego	string
	author_avatar	avatar autora	uri
	author_avatar_med	avatar autora (średni rozmiar)	uri
	author_avatar_lo	avatar autora (mały rozmiar)	uri
	author_group	grupa autora	int
	preview	miniatura	uri
	user_lists	listy ulubionych na których znajduje się link	int[]
	plus18	link dla dorosłych	bool
	status
	can_vote	czy uzytkownik może głosować	bool
	has_own_content	czy link ma treść	bool
	category	domena grupy	string
	Poniżesz pola zostaną wypełnione jeśli do pytania o link lub listę linków zostanie dołożony userkey
	user_vote	dig jeśli użytkownik wykopał ten link, lub bury, jeśli zakopał	string
	user_observe	true, jeśli użytkownik obserwuje ten link	bool
	user_favorite	true, jeśli użytkownik dodał link do ulubionych	bool
*/

type Link struct {
	ID           int `json:"id"`
	Title        string
	Description  string
	Tags         string
	Url          string
	SourceUrl    string
	VoteCount    int
	CommentCount int
	ReportCount  int
	Date         string

	AuthorInfo

	Type          string
	Group         string
	Preview       string
	UserLists     []int
	Plus18        bool
	Status        string
	CanVote       bool
	IsHot         bool
	HasOwnContent bool
	Category      string
	CategoryName  string
	//
	UserVote     bool
	UserObserve  bool
	UserFavorite bool
	ViolationUrl string
	Info         string
	App          string
	OwnContent   string
}

type UserGroup int

const (
	GroupGreen   UserGroup = 0
	GroupOrange            = 1
	GroupBrown             = 2
	GroupAdmin             = 5
	GroupBanned            = 1001
	GroupDeleted           = 1002
	GroupClient            = 2001
)

func (u UserGroup) Color() string {
	switch u {
	case GroupGreen:
		return "#339933"
	case GroupOrange:
		return "#339933"
	case GroupBrown:
		return "#BB0000"
	case GroupAdmin:
		return "#000000"
	case GroupBanned:
		fallthrough
	case GroupDeleted:
		return "#999999"
	case GroupClient:
		return "#3F6FA0"
	}
	panic("unknown group")
}

func (u UserGroup) Name() string {
	switch u {
	case GroupGreen:
		return "Zielony"
	case GroupOrange:
		return "Pomarańczowy"
	case GroupBrown:
		return "Bordowy"
	case GroupAdmin:
		return "Administrator"
	case GroupBanned:
		return "Zbanowany"
	case GroupDeleted:
		return "Usunięty"
	case GroupClient:
		return "Klient"
	}
	panic("unknown group")
}

type WrappedError struct {
	Error Error
}

type Error struct {
	Code    int
	Message string
}

/*
BURY
Pole	Wartość	Typ
reason	identyfikator powodu zakopu	int
author	autor	string
author_avatar	avatar autora	uri
author_avatar_med	avatar autora (średni rozmiar)	uri
author_avatar_lo	avatar autora (mały rozmiar)	uri
author_group	grupa autora	int
*/
type Bury struct {
	Reason int
	AuthorInfo
}

type AuthorInfo struct {
	Author          string
	AuthorAvatar    string
	AuthorAvatarBig string
	AuthorAvatarMed string
	AuthorAvatarLo  string
	AuthorGroup     string
	AuthorSex       string
}

/*
DIG
Pole	Wartość	Typ
author	autor	string
author_avatar	avatar autora	uri
author_avatar_med	avatar autora (średni rozmiar)	uri
author_avatar_lo	avatar autora (mały rozmiar)	uri
author_group	grupa autora	int
*/
type Dig struct {
	AuthorInfo
}

/*
COMMENT
Pole	Wartość	Typ
id	identyfikator komentarza	int
date	data komentarza	date
author	autor komentarza	string
author_avatar	avatar autora	uri
author_avatar_med	avatar autora (średni rozmiar)	uri
author_avatar_lo	avatar autora (mały rozmiar)	uri
author_group	grupa autora	int
vote_count	liczba głosów	int
body	treść	string
parent_id	identyfikator komentarza nadrzędnego	int
status	Status komentarza (own/new/readed)	string
embed	obrazek lub film dołączony do obiektu	embed
link	link do którego dodany jest komentarz (występuje tylko w metodzie profile/comments)	Link
*/
type Comment struct {
	ID   int
	Date string

	AuthorInfo

	VoteCount int
	Body      string
	ParentID  int
	Status    string
	Embed     string
	Link      string
}

/*
RELATEDLINK
Pole	Wartość	Typ
id	identyfikator linku	int
url	adres URL linku	uri
title	tytuł linku	string
plus18	link dla dorosłych	bool
vote_count	ilość głosów	int
entry_count	ilość wejść z linka (dla linków trackback)	int
user_vote	głos zalogowanego użytkownika (+1 / -1 / null jeśli brak głosu)	int
author	login dodającego	string
author_avatar	avatar autora	uri
author_avatar_med	avatar autora (średni rozmiar)	uri
author_avatar_lo	avatar autora (mały rozmiar)	uri
author_group	grupa autora	int
link	link do którego dodany jest link powiązany (występuje tylko w metodzie profile/related)	Link
*/
type RelatedLink struct {
	ID         int
	Url        string
	Title      string
	Plus18     bool
	VoteCount  int
	EntryCount int
	UserVote   int

	AuthorInfo

	Link string
}

type LinkResource interface {
	Index(id int) (*Link, error)
	//Dig()
	//Cancel()
	//Bury()
	Comments(id int) ([]*Comment, error)
	Reports(id int) ([]*Bury, error)
	Digs(id int) ([]*Dig, error)
	Related(id int) ([]*RelatedLink, error)
	BuryReasons()
	//Observe()
	//Favorite()
}
