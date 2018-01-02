package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/coreos/pkg/httputil"
	"github.com/stretchr/testify/assert"
)

type fakeRT struct {
	call func(*http.Request) (*http.Response, error)
}

func (f *fakeRT) RoundTrip(r *http.Request) (*http.Response, error) {
	return f.call(r)
}

func TestNew(t *testing.T) {
	call := false
	c := New("app", "user", &fakeRT{func(r *http.Request) (*http.Response, error) {
		call = true
		rec := httptest.NewRecorder()
		//httputil.WriteJSONResponse(rec, 200, )
		return rec.Result(), nil
	}})
	req := &Request{RequestMethod: http.MethodGet, Resource: "links", Method: "promoted"}
	assert.NoError(t, c.DoAndParse(req, nil))
	assert.True(t, call)
}

func TestLinksPromoted(t *testing.T) {
	c := New("app", "user", &fakeRT{func(r *http.Request) (*http.Response, error) {
		rec := httptest.NewRecorder()
		links := []*Link{{ID: 10}}
		httputil.WriteJSONResponse(rec, 200, links)
		return rec.Result(), nil
	}})
	links, err := c.Links().Promoted(1, PromotedByDay)
	assert.NoError(t, err)
	assert.NotEmpty(t, links)
}

func TestLinksUpcoming(t *testing.T) {
	c := New("app", "user", &fakeRT{func(r *http.Request) (*http.Response, error) {
		rec := httptest.NewRecorder()
		links := []*Link{{ID: 10}}
		httputil.WriteJSONResponse(rec, 200, links)
		return rec.Result(), nil
	}})
	links, err := c.Links().Upcoming(1, UpcomingVotes)
	assert.NoError(t, err)
	assert.NotEmpty(t, links)
}

func TestLinkIndex(t *testing.T) {
	c := New("app", "user", &fakeRT{func(r *http.Request) (*http.Response, error) {
		rec := httptest.NewRecorder()
		link := &Link{ID: 10}
		httputil.WriteJSONResponse(rec, 200, link)
		return rec.Result(), nil
	}})
	link, err := c.Link().Index(10)
	assert.NoError(t, err)
	assert.NotNil(t, link)
}

func TestLinkIndex2(t *testing.T) {
	c := New("app", "user", &fakeRT{func(r *http.Request) (*http.Response, error) {
		rec := httptest.NewRecorder()
		rec.Header().Set("Content-Type", httputil.JSONContentType)
		f, err := os.Open("testdata/link2.json")
		if err != nil {
			return nil, err
		}
		defer f.Close()
		io.Copy(rec, f)
		return rec.Result(), nil
	}})
	link, err := c.Link().Index(1475611)
	assert.NoError(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, 1475611, link.ID)
}

func TestLinkDigs(t *testing.T) {
	c := New("app", "user", &fakeRT{func(r *http.Request) (*http.Response, error) {
		rec := httptest.NewRecorder()
		link := []*Dig{{AuthorInfo{Author: "Imię"}}}
		httputil.WriteJSONResponse(rec, 200, link)
		return rec.Result(), nil
	}})
	link, err := c.Link().Digs(10)
	assert.NoError(t, err)
	assert.NotNil(t, link)
}

func TestLinkReports(t *testing.T) {

}
