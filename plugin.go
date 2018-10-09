package main

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"io"
	"time"
)

type (
	Config struct {
		Loadbalancer string
		Tag          string
		User         string
		KeyPath      string
		Key          string
		Password     string
		Port         int
		Timeout      time.Duration
		Pat          string
		SourcePath   string
		TargetPath   string
		PreSync      []string
		PostSync     []string
	}

	Plugin struct {
		Config Config
		Hosts  []Host
		Writer io.Writer
	}

	Host struct {
		ID int
		IP string
	}

	TokenSource struct {
		AccessToken string
	}
)

const (
	onlyOneHost          = "please provide only a loadbalancer or a tag"
	noHost               = "please provide a loadbalancer or a tag"
	tagNotFound          = "provided tag not found or no droplets assigned"
	loadbalancerNotFound = "provided loadbalancer not found"
	onlyOneCredential    = "please provide a key or a password"
	noCredentials        = "no key or password provided"
	noPAT                = "no PAT provided"
	noDropletsFound      = "no droplets found"
	couldntResolveIPs    = "could not resolve IDs to IPs"
)

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// Validate checks the values of the config parameters.
// If the examined parameter did not pass the check, either a default value
// is set or, if this is not possible or useful, an error is returned.
func (c *Config) Validate() error {
	// asset both set
	if len(c.Loadbalancer) > 0 && len(c.Tag) > 0 {
		return fmt.Errorf(onlyOneHost)
	}

	// assert none set
	if len(c.Loadbalancer) == 0 && len(c.Tag) == 0 {
		return fmt.Errorf(noHost)
	}

	// asset user unset
	if len(c.User) == 0 {
		c.User = "root"
	}

	// assert more than one set
	if (len(c.KeyPath) > 0 || len(c.Key) > 0) && len(c.Password) > 0 {
		return fmt.Errorf(onlyOneCredential)
	}

	// assert none set
	if len(c.KeyPath) == 0 && len(c.Key) == 0 && len(c.Password) == 0 {
		return fmt.Errorf(noCredentials)
	}

	// assert port unset
	if c.Port <= 0 {
		c.Port = 22
	}

	// assert timeout unset
	if c.Timeout == time.Duration(0) {
		c.Timeout = time.Duration(60)
	}

	// assert token unset
	if len(c.Pat) == 0 {
		return fmt.Errorf(noPAT)
	}

	// assert source unset
	if len(c.SourcePath) == 0 {
		c.SourcePath = "./"
	}

	// assert target unset
	if len(c.TargetPath) == 0 {
		c.TargetPath = "./"
	}

	return nil
}

func (p Plugin) getClient() *godo.Client {
	tokenSource := &TokenSource{
		AccessToken: p.Config.Pat,
	}

	oauthClient := oauth2.NewClient(context.TODO(), tokenSource)
	return godo.NewClient(oauthClient)
}

// getDropletIDsByLoadbalancer searches for a loadbalancer with the
// specified name and, if found, returns its droplets as Hosts.
func (p Plugin) getHostsByLoadbalancer() ([]Host, error) {
	client := p.getClient()

	// get all loadbalancer
	lbs, _, err := client.LoadBalancers.List(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	// get all droplets
	droplets, _, err := client.Droplets.List(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	// intersect loadbalancer
	var ids []int
	for _, lb := range lbs {
		if lb.Name == p.Config.Loadbalancer {
			ids = lb.DropletIDs
		}
	}

	// intersect loadbalancer and droplets
	var hosts []Host
	for _, droplet := range droplets {
		for _, id := range ids {
			if droplet.ID == id {
				if ip, err := droplet.PublicIPv4(); err == nil {
					hosts = append(hosts, Host{ID: id, IP: ip})
				}
			}
		}
	}

	if len(hosts) == 0 {
		return nil, fmt.Errorf(noDropletsFound)
	}

	return hosts, nil
}

// getHostsByTag searches for a tag with the specified name and, if
// found, returns its droplets as Hosts.
func (p Plugin) getHostsByTag() ([]Host, error) {
	client := p.getClient()

	// get droplets with tag
	droplets, _, err := client.Droplets.ListByTag(context.TODO(), p.Config.Tag, nil)
	if err != nil {
		return nil, err
	}

	// covert to Droplet to Host
	var hosts []Host
	for _, droplet := range droplets {
		if ip, err := droplet.PublicIPv4(); err == nil {
			hosts = append(hosts, Host{ID: droplet.ID, IP: ip})
		}
	}

	if len(hosts) == 0 {
		return nil, fmt.Errorf(noDropletsFound)
	}

	return hosts, nil
}

func (p Plugin) runPreSyncScript() error {


	return nil
}

func (p *Plugin) Exec() error {
	// check command line arguments
	if err := p.Config.Validate(); err != nil {
		return err
	}

	if len(p.Config.Loadbalancer) != 0 {
		hosts, err := p.getHostsByLoadbalancer()
		if err != nil {
			return err
		}
		p.Hosts = hosts
	}

	if len(p.Config.Tag) != 0 {
		hosts, err := p.getHostsByTag()
		if err != nil {
			return err
		}
		p.Hosts = hosts
	}

	fmt.Println(p.Hosts)

	fmt.Println("pre sync", p.Config.PreSync)

	// TODO: pre sync script

	// TODO: sync

	// TODO: post sync script

	return nil
}