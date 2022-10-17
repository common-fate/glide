package ssov2

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	organizationTypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"golang.org/x/sync/errgroup"
)

type OrganizationGraph struct {
	Root    *Node
	idMap   map[string]*Node
	idMapMu sync.Mutex
}
type Node struct {
	descendantsMu sync.Mutex
	ID            string
	Graph         *OrganizationGraph
	Parent        *Node
	// Direct children of this node
	Children []*Node
	// All descendants of this node
	Descendants        []*Node
	OrganizationalUnit *organizationTypes.OrganizationalUnit
	Account            *organizationTypes.Account
	Root               *organizationTypes.Root
}

func (p *Provider) buildOrganizationGraph(ctx context.Context) (*OrganizationGraph, error) {
	roots, err := p.orgClient.ListRoots(ctx, &organizations.ListRootsInput{})
	if err != nil {
		return nil, err
	}

	if len(roots.Roots) != 1 {
		return nil, fmt.Errorf("expected to find 1 organization root but found %v", len(roots.Roots))
	}
	root := roots.Roots[0]
	graph := OrganizationGraph{
		Root: &Node{
			ID:   aws.ToString(root.Id),
			Root: &root,
		},
	}

	graph.Root.Graph = &graph

	graph.idMap = map[string]*Node{aws.ToString(root.Id): graph.Root}
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
			accountIDs = append(accountIDs, aws.ToString(child.Account.Id))
		}
	}
	return accountIDs
}
func (n *Node) DescendantOrganizationalUnitIDs() []string {
	var accountIDs []string
	for _, child := range n.Descendants {
		if child.IsOrganizationalUnit() {
			accountIDs = append(accountIDs, aws.ToString(child.OrganizationalUnit.Id))
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
		g, gctx := errgroup.WithContext(ctx)
		childOUs, err := provider.listChildOusForParent(ctx, n.ID)
		if err != nil {
			return err
		}
		for i := range childOUs {
			node := &Node{
				ID:                 aws.ToString(childOUs[i].Id),
				OrganizationalUnit: &childOUs[i],
				Parent:             n,
				Graph:              n.Graph,
			}

			n.Children = append(n.Children, node)
			n.Graph.idMapMu.Lock()
			n.Graph.idMap[node.ID] = node
			n.Graph.idMapMu.Unlock()
			g.Go(func() error {
				return node.BuildGraph(gctx, provider)
			})

		}
		g.Go(func() error {
			childAccounts, err := provider.listChildAccountsForParent(ctx, n.ID)
			if err != nil {
				return err
			}
			for i := range childAccounts {
				node := &Node{
					ID:      aws.ToString(childAccounts[i].Id),
					Account: &childAccounts[i],
					Parent:  n,
					Graph:   n.Graph,
				}
				n.Children = append(n.Children, node)
				n.Graph.idMapMu.Lock()
				n.Graph.idMap[node.ID] = node
				n.Graph.idMapMu.Unlock()
			}
			return nil
		})

		err = g.Wait()
		if err != nil {
			return err
		}
		// assign all children as decendants
		n.Descendants = append(n.Descendants, n.Children...)
		if n.Parent != nil {
			// append decendants to parent
			n.Parent.descendantsMu.Lock()
			n.Parent.Descendants = append(n.Parent.Descendants, n.Descendants...)
			n.Parent.descendantsMu.Unlock()
		}

	}
	return nil
}
