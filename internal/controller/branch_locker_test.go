package controller

import "testing"

func TestBranchLocker_TryAcquireAndRelease(t *testing.T) {
	bl := NewBranchLocker()

	ok, holder := bl.TryAcquire("ws:feature-1", "task-a")
	if !ok {
		t.Fatalf("Expected TryAcquire to succeed, got held by %q", holder)
	}

	// Same task re-acquiring is idempotent.
	ok, _ = bl.TryAcquire("ws:feature-1", "task-a")
	if !ok {
		t.Fatal("Expected idempotent TryAcquire to succeed")
	}

	// Different task should be rejected.
	ok, holder = bl.TryAcquire("ws:feature-1", "task-b")
	if ok {
		t.Fatal("Expected TryAcquire to fail for different task")
	}
	if holder != "task-a" {
		t.Errorf("Expected holder %q, got %q", "task-a", holder)
	}

	// Release and re-acquire by different task.
	bl.Release("ws:feature-1", "task-a")
	ok, _ = bl.TryAcquire("ws:feature-1", "task-b")
	if !ok {
		t.Fatal("Expected TryAcquire to succeed after release")
	}
}

func TestBranchLocker_DifferentKeysIndependent(t *testing.T) {
	bl := NewBranchLocker()

	ok, _ := bl.TryAcquire("ws-a:feature-1", "task-a")
	if !ok {
		t.Fatal("Expected TryAcquire to succeed")
	}

	// Different key should succeed independently.
	ok, _ = bl.TryAcquire("ws-b:feature-1", "task-b")
	if !ok {
		t.Fatal("Expected TryAcquire on different key to succeed")
	}
}

func TestBranchLocker_ReleaseWrongTask(t *testing.T) {
	bl := NewBranchLocker()

	bl.TryAcquire("ws:feature-1", "task-a")

	// Releasing with wrong task name should be a no-op.
	bl.Release("ws:feature-1", "task-b")

	if h := bl.Holder("ws:feature-1"); h != "task-a" {
		t.Errorf("Expected holder %q after wrong release, got %q", "task-a", h)
	}
}

func TestBranchLocker_ReleaseUnheldKey(t *testing.T) {
	bl := NewBranchLocker()

	// Should not panic.
	bl.Release("ws:nonexistent", "task-a")
}

func TestBranchLocker_Holder(t *testing.T) {
	bl := NewBranchLocker()

	if h := bl.Holder("ws:feature-1"); h != "" {
		t.Errorf("Expected empty holder, got %q", h)
	}

	bl.TryAcquire("ws:feature-1", "task-a")
	if h := bl.Holder("ws:feature-1"); h != "task-a" {
		t.Errorf("Expected holder %q, got %q", "task-a", h)
	}
}
