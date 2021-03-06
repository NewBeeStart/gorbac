package gorbac

import "fmt"

// InherCircle returns an error when detecting any circle inheritance.
func InherCircle(rbac *RBAC) (err error) {
	rbac.mutex.Lock()

	skipped := make(map[string]struct{})
	var stack []string

	for id := range rbac.roles {
		if err = dfs(rbac, id, skipped, stack); err != nil {
			break
		}
	}
	rbac.mutex.Unlock()
	return err
}

func dfs(rbac *RBAC, id string, skipped map[string]struct{}, stack []string) error {
	if _, ok := skipped[id]; ok {
		return nil
	}
	for _, item := range stack {
		if item == id {
			return fmt.Errorf("Found circle: %s", stack)
		}
	}
	if len(rbac.parents[id]) == 0 {
		stack = make([]string, 0)
		skipped[id] = empty
		return nil
	}
	stack = append(stack, id)
	for pid := range rbac.parents[id] {
		if err := dfs(rbac, pid, skipped, stack); err != nil {
			return err
		}
	}
	return nil
}

// AnyGranted checks if any role has the permission.
func AnyGranted(rbac *RBAC, roles []string, permission Permission,
	assert AssertionFunc) (rslt bool) {
	rbac.mutex.Lock()
	for _, role := range roles {
		if rbac.isGranted(role, permission, assert) {
			rslt = true
			break
		}
	}
	rbac.mutex.Unlock()
	return rslt
}

// AllGranted checks if all roles have the permission.
func AllGranted(rbac *RBAC, roles []string, permission Permission,
	assert AssertionFunc) (rslt bool) {
	rbac.mutex.Lock()
	for _, role := range roles {
		if !rbac.isGranted(role, permission, assert) {
			rslt = true
			break
		}
	}
	rbac.mutex.Unlock()
	return !rslt
}
