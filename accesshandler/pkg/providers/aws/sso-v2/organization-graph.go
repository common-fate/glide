package ssov2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	organizationTypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

type Node struct {
	ID     string
	Graph  *OrganizationGraph
	Parent *Node
	// Direct children of this node
	Children []*Node
	// All descendants of this node
	Descendants        []*Node
	OrganizationalUnit *organizationTypes.OrganizationalUnit
	Account            *organizationTypes.Account
	Root               *organizationTypes.Root
}
type OrganizationGraph struct {
	Root  *Node
	idMap map[string]*Node
}

func (p *Provider) buildOrganizationGraph(ctx context.Context) (*OrganizationGraph, error) {
	roots, err := p.orgClient.ListRoots(ctx, &organizations.ListRootsInput{})
	if err != nil {
		return nil, err
	}
	graph := OrganizationGraph{
		Root: &Node{
			ID:   *roots.Roots[0].Id,
			Root: &roots.Roots[0],
		},
	}
	graph.Root.Graph = &graph

	graph.idMap = map[string]*Node{*roots.Roots[0].Id: graph.Root}
	err = graph.Root.BuildGraph(ctx, p)
	if err != nil {
		return nil, err
	}
	return &graph, nil
}

func (n *Node) IsRoot() bool {
	return n.Root != nil
}
func (n *Node) IsAccount() bool {
	return n.Account != nil
}
func (n *Node) IsOrganizationalUnit() bool {
	return n.OrganizationalUnit != nil
}

func (n *Node) DescendantAccountIDs() []string {
	var accountIDs []string
	for _, child := range n.Descendants {
		if child.IsAccount() {
			accountIDs = append(accountIDs, *child.Account.Id)
		}
	}
	return accountIDs
}
func (n *Node) DescendantOrganizationalUnitIDs() []string {
	var accountIDs []string
	for _, child := range n.Descendants {
		if child.IsOrganizationalUnit() {
			accountIDs = append(accountIDs, *child.OrganizationalUnit.Id)
		}
	}
	return accountIDs
}
func (n *Node) DescendantOrganisationTypeAccounts() []organizationTypes.Account {
	var accounts []organizationTypes.Account
	for i := range n.Descendants {
		if n.Descendants[i].IsAccount() {
			accounts = append(accounts, *n.Descendants[i].Account)
		}
	}
	return accounts
}
func (n *Node) DescendantAccounts() []*Node {
	var accounts []*Node
	for i := range n.Descendants {
		if n.Descendants[i].IsAccount() {
			accounts = append(accounts, n.Descendants[i])
		}
	}
	return accounts
}
func (n *Node) DescendantOrganizationalUnits() []*Node {
	var organizationUnits []*Node
	for i := range n.Descendants {
		if n.Descendants[i].IsOrganizationalUnit() {
			organizationUnits = append(organizationUnits, n.Descendants[i])
		}
	}
	return organizationUnits
}

func (n *Node) BuildGraph(ctx context.Context, provider *Provider) error {
	if n.IsOrganizationalUnit() || n.IsRoot() {
		childOUs, err := provider.listChildOusForParent(ctx, n.ID)
		if err != nil {
			return err
		}
		for i := range childOUs {
			node := &Node{
				ID:                 *childOUs[i].Id,
				OrganizationalUnit: &childOUs[i],
				Parent:             n,
				Graph:              n.Graph,
			}
			n.Children = append(n.Children, node)
			n.Graph.idMap[node.ID] = node
		}
		childAccounts, err := provider.listChildAccountsForParent(ctx, n.ID)
		if err != nil {
			return err
		}
		for i := range childAccounts {
			node := &Node{
				ID:      *childAccounts[i].Id,
				Account: &childAccounts[i],
				Parent:  n,
				Graph:   n.Graph,
			}
			n.Children = append(n.Children, node)
			n.Graph.idMap[node.ID] = node
		}
		// assign all children as decendants
		n.Descendants = n.Children
		for i := range n.Children {
			err = n.Children[i].BuildGraph(ctx, provider)
			if err != nil {
				return err
			}
		}
		if n.Parent != nil {
			// append decendants to parent
			n.Parent.Descendants = append(n.Parent.Descendants, n.Descendants...)
		}

	}
	return nil
}
