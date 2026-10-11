package rdb_test

import (
	"context"
	"errors"
	"net/netip"
	"os"
	"testing"

	"exp-myddns/internal/generated/rdb"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultDatabaseURL = "postgres://postgres:postgres@localhost:5432/mydatabase?sslmode=disable"

	testUserID        = "01890000-0000-7000-8000-000000000001"
	testOtherUserID   = "01890000-0000-7000-8000-000000000002"
	testDeviceID      = "01890000-0000-7000-8000-000000000011"
	testOtherDeviceID = "01890000-0000-7000-8000-000000000012"
	testRecordID      = "01890000-0000-7000-8000-000000000021"
	testOtherRecordID = "01890000-0000-7000-8000-000000000022"
)

var testPool *pgxpool.Pool

var cmpOptions = cmp.Options{
	cmpopts.EquateComparable(netip.Addr{}),
	cmpopts.IgnoreFields(rdb.User{}, "CreatedAt", "UpdatedAt"),
	cmpopts.IgnoreFields(rdb.Device{}, "CreatedAt"),
	cmpopts.IgnoreFields(rdb.Record{}, "ReportedAt"),
	cmpopts.IgnoreFields(rdb.Credential{}, "CreatedAt", "UpdatedAt"),
}

func TestMain(m *testing.M) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = defaultDatabaseURL
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic("failed to create database pool: " + err.Error())
	}
	testPool = pool

	code := m.Run()
	pool.Close()
	os.Exit(code)
}

func newTestQueries(t *testing.T) (*rdb.Queries, context.Context) {
	t.Helper()

	ctx := t.Context()
	tx, err := testPool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	t.Cleanup(func() {
		if err := tx.Rollback(context.Background()); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			t.Errorf("failed to rollback transaction: %v", err)
		}
	})

	return rdb.New(tx), ctx
}

func mustUUID(t *testing.T, s string) pgtype.UUID {
	t.Helper()

	var u pgtype.UUID
	if err := u.Scan(s); err != nil {
		t.Fatalf("failed to parse uuid %q: %v", s, err)
	}

	return u
}

func mustAddr(t *testing.T, s string) *netip.Addr {
	t.Helper()

	addr := netip.MustParseAddr(s)

	return &addr
}

func seedUser(t *testing.T, q *rdb.Queries, ctx context.Context, id, userName, displayName string) rdb.User {
	t.Helper()

	u, err := q.CreateUser(ctx, rdb.CreateUserParams{
		ID:          mustUUID(t, id),
		UserName:    userName,
		DisplayName: displayName,
	})
	if err != nil {
		t.Fatalf("seed user: CreateUser() error = %v", err)
	}

	return u
}

func seedDevice(t *testing.T, q *rdb.Queries, ctx context.Context, id, userID, name, memo string) rdb.Device {
	t.Helper()

	d, err := q.CreateDevice(ctx, rdb.CreateDeviceParams{
		ID:     mustUUID(t, id),
		UserID: mustUUID(t, userID),
		Name:   name,
		Memo:   memo,
	})
	if err != nil {
		t.Fatalf("seed device: CreateDevice() error = %v", err)
	}

	return d
}

func seedRecord(t *testing.T, q *rdb.Queries, ctx context.Context, id, deviceID, recordType, ipAddress string) rdb.Record {
	t.Helper()

	r, err := q.CreateRecord(ctx, rdb.CreateRecordParams{
		ID:        mustUUID(t, id),
		DeviceID:  mustUUID(t, deviceID),
		Type:      recordType,
		IpAddress: mustAddr(t, ipAddress),
	})
	if err != nil {
		t.Fatalf("seed record: CreateRecord() error = %v", err)
	}

	return r
}

func seedCredential(t *testing.T, q *rdb.Queries, ctx context.Context, userID, provider, secretHash, subject string) rdb.Credential {
	t.Helper()

	c, err := q.CreateCredential(ctx, rdb.CreateCredentialParams{
		UserID:     mustUUID(t, userID),
		Provider:   provider,
		SecretHash: pgtype.Text{String: secretHash, Valid: secretHash != ""},
		Subject:    subject,
	})
	if err != nil {
		t.Fatalf("seed credential: CreateCredential() error = %v", err)
	}

	return c
}

func TestGetUserByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userId  string
		want    rdb.User
		wantErr error
	}{
		{
			name:   "success",
			userId: testUserID,
			want: rdb.User{
				ID:          mustUUID(t, testUserID),
				UserName:    "alice",
				DisplayName: "Alice",
			},
		},
		{
			name:    "not found",
			userId:  testOtherUserID,
			wantErr: pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")

			got, err := q.GetUserByID(ctx, mustUUID(t, tt.userId))
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetUserByID() error = %v", err)
			}
			if diff := cmp.Diff(tt.want, got, cmpOptions); diff != "" {
				t.Errorf("GetUserByID() mismatch (-want +got):\n%s", diff)
			}
			if !got.CreatedAt.Valid || !got.UpdatedAt.Valid {
				t.Errorf("GetUserByID() timestamps are invalid: %+v", got)
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		seed    bool
		wantLen int
	}{
		{
			name:    "empty",
			seed:    false,
			wantLen: 0,
		},
		{
			name:    "success",
			seed:    true,
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			if tt.seed {
				seedUser(t, q, ctx, testUserID, "alice", "Alice")
				seedUser(t, q, ctx, testOtherUserID, "bob", "Bob")
			}

			got, err := q.ListUsers(ctx)
			if err != nil {
				t.Fatalf("ListUsers() error = %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("ListUsers() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		userName    string
		displayName string
	}{
		{
			name:        "success",
			userName:    "alice",
			displayName: "Alice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)

			got, err := q.CreateUser(ctx, rdb.CreateUserParams{
				ID:          mustUUID(t, testUserID),
				UserName:    tt.userName,
				DisplayName: tt.displayName,
			})
			if err != nil {
				t.Fatalf("CreateUser() error = %v", err)
			}

			want := rdb.User{
				ID:          mustUUID(t, testUserID),
				UserName:    tt.userName,
				DisplayName: tt.displayName,
			}
			if diff := cmp.Diff(want, got, cmpOptions); diff != "" {
				t.Errorf("CreateUser() mismatch (-want +got):\n%s", diff)
			}
			if !got.CreatedAt.Valid || !got.UpdatedAt.Valid {
				t.Errorf("CreateUser() timestamps are invalid: %+v", got)
			}

			found, err := q.GetUserByID(ctx, got.ID)
			if err != nil {
				t.Fatalf("GetUserByID() after CreateUser() error = %v", err)
			}
			if diff := cmp.Diff(want, found, cmpOptions); diff != "" {
				t.Errorf("GetUserByID() after CreateUser() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		userID      string
		userName    string
		displayName string
		wantErr     error
	}{
		{
			name:        "success",
			userID:      testUserID,
			userName:    "alice-updated",
			displayName: "Alice Updated",
		},
		{
			name:        "not found",
			userID:      testOtherUserID,
			userName:    "nobody",
			displayName: "Nobody",
			wantErr:     pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			created := seedUser(t, q, ctx, testUserID, "alice", "Alice")

			got, err := q.UpdateUser(ctx, rdb.UpdateUserParams{
				ID:          mustUUID(t, tt.userID),
				UserName:    tt.userName,
				DisplayName: tt.displayName,
			})
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("UpdateUser() error = %v", err)
			}

			want := rdb.User{
				ID:          mustUUID(t, tt.userID),
				UserName:    tt.userName,
				DisplayName: tt.displayName,
			}
			if diff := cmp.Diff(want, got, cmpOptions); diff != "" {
				t.Errorf("UpdateUser() mismatch (-want +got):\n%s", diff)
			}
			if got.CreatedAt.Time != created.CreatedAt.Time {
				t.Errorf("UpdateUser() changed created_at: got %v, want %v", got.CreatedAt.Time, created.CreatedAt.Time)
			}
			if !got.UpdatedAt.Valid || got.UpdatedAt.Time.Before(created.UpdatedAt.Time) {
				t.Errorf("UpdateUser() updated_at = %v, want not before %v", got.UpdatedAt.Time, created.UpdatedAt.Time)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  string
		wantErr error
	}{
		{
			name:   "success",
			userID: testUserID,
		},
		{
			name:   "not found",
			userID: testOtherUserID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")

			if err := q.DeleteUser(ctx, mustUUID(t, tt.userID)); err != nil {
				t.Fatalf("DeleteUser() error = %v", err)
			}

			_, err := q.GetUserByID(ctx, mustUUID(t, tt.userID))
			if !errors.Is(err, pgx.ErrNoRows) {
				t.Fatalf("GetUserByID() after DeleteUser() error = %v, want %v", err, pgx.ErrNoRows)
			}
		})
	}
}

func TestGetDeviceByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		deviceID string
		want     rdb.Device
		wantErr  error
	}{
		{
			name:     "success",
			deviceID: testDeviceID,
			want: rdb.Device{
				ID:     mustUUID(t, testDeviceID),
				UserID: mustUUID(t, testUserID),
				Name:   "router",
				Memo:   "living room",
			},
		},
		{
			name:     "not found",
			deviceID: testOtherDeviceID,
			wantErr:  pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")
			seedDevice(t, q, ctx, testDeviceID, testUserID, "router", "living room")

			got, err := q.GetDeviceByID(ctx, mustUUID(t, tt.deviceID))
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("GetDeviceByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetDeviceByID() error = %v", err)
			}
			if diff := cmp.Diff(tt.want, got, cmpOptions); diff != "" {
				t.Errorf("GetDeviceByID() mismatch (-want +got):\n%s", diff)
			}
			if !got.CreatedAt.Valid {
				t.Errorf("GetDeviceByID() created_at is invalid: %+v", got)
			}
		})
	}
}

func TestListDevicesByUserID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  string
		seed    bool
		wantLen int
	}{
		{
			name:    "empty",
			userID:  testUserID,
			seed:    false,
			wantLen: 0,
		},
		{
			name:    "success",
			userID:  testUserID,
			seed:    true,
			wantLen: 2,
		},
		{
			name:    "filtered by other user",
			userID:  testOtherUserID,
			seed:    true,
			wantLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			if tt.seed {
				seedUser(t, q, ctx, testUserID, "alice", "Alice")
				seedDevice(t, q, ctx, testDeviceID, testUserID, "router", "living room")
				seedDevice(t, q, ctx, testOtherDeviceID, testUserID, "nas", "study")
			}

			got, err := q.ListDevicesByUserID(ctx, mustUUID(t, tt.userID))
			if err != nil {
				t.Fatalf("ListDevicesByUserID() error = %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("ListDevicesByUserID() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestCreateDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		deviceName string
		memo       string
	}{
		{
			name:       "success",
			deviceName: "router",
			memo:       "living room",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")

			got, err := q.CreateDevice(ctx, rdb.CreateDeviceParams{
				ID:     mustUUID(t, testDeviceID),
				UserID: mustUUID(t, testUserID),
				Name:   tt.deviceName,
				Memo:   tt.memo,
			})
			if err != nil {
				t.Fatalf("CreateDevice() error = %v", err)
			}

			want := rdb.Device{
				ID:     mustUUID(t, testDeviceID),
				UserID: mustUUID(t, testUserID),
				Name:   tt.deviceName,
				Memo:   tt.memo,
			}
			if diff := cmp.Diff(want, got, cmpOptions); diff != "" {
				t.Errorf("CreateDevice() mismatch (-want +got):\n%s", diff)
			}
			if !got.CreatedAt.Valid {
				t.Errorf("CreateDevice() created_at is invalid: %+v", got)
			}

			found, err := q.GetDeviceByID(ctx, got.ID)
			if err != nil {
				t.Fatalf("GetDeviceByID() after CreateDevice() error = %v", err)
			}
			if diff := cmp.Diff(want, found, cmpOptions); diff != "" {
				t.Errorf("GetDeviceByID() after CreateDevice() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		deviceID   string
		deviceName string
		memo       string
		wantErr    error
	}{
		{
			name:       "success",
			deviceID:   testDeviceID,
			deviceName: "router-updated",
			memo:       "bedroom",
		},
		{
			name:       "not found",
			deviceID:   testOtherDeviceID,
			deviceName: "ghost",
			memo:       "nowhere",
			wantErr:    pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")
			created := seedDevice(t, q, ctx, testDeviceID, testUserID, "router", "living room")

			got, err := q.UpdateDevice(ctx, rdb.UpdateDeviceParams{
				ID:   mustUUID(t, tt.deviceID),
				Name: tt.deviceName,
				Memo: tt.memo,
			})
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("UpdateDevice() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("UpdateDevice() error = %v", err)
			}

			want := rdb.Device{
				ID:     mustUUID(t, tt.deviceID),
				UserID: mustUUID(t, testUserID),
				Name:   tt.deviceName,
				Memo:   tt.memo,
			}
			if diff := cmp.Diff(want, got, cmpOptions); diff != "" {
				t.Errorf("UpdateDevice() mismatch (-want +got):\n%s", diff)
			}
			if got.CreatedAt.Time != created.CreatedAt.Time {
				t.Errorf("UpdateDevice() changed created_at: got %v, want %v", got.CreatedAt.Time, created.CreatedAt.Time)
			}
		})
	}
}

func TestDeleteDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		deviceID string
	}{
		{
			name:     "success",
			deviceID: testDeviceID,
		},
		{
			name:     "not found",
			deviceID: testOtherDeviceID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")
			seedDevice(t, q, ctx, testDeviceID, testUserID, "router", "living room")

			if err := q.DeleteDevice(ctx, mustUUID(t, tt.deviceID)); err != nil {
				t.Fatalf("DeleteDevice() error = %v", err)
			}

			_, err := q.GetDeviceByID(ctx, mustUUID(t, tt.deviceID))
			if !errors.Is(err, pgx.ErrNoRows) {
				t.Fatalf("GetDeviceByID() after DeleteDevice() error = %v, want %v", err, pgx.ErrNoRows)
			}
		})
	}
}

func TestGetRecordByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		recordID string
		want     rdb.Record
		wantErr  error
	}{
		{
			name:     "success",
			recordID: testRecordID,
			want: rdb.Record{
				ID:        mustUUID(t, testRecordID),
				DeviceID:  mustUUID(t, testDeviceID),
				Type:      "A",
				IpAddress: mustAddr(t, "192.0.2.1"),
			},
		},
		{
			name:     "not found",
			recordID: testOtherRecordID,
			wantErr:  pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")
			seedDevice(t, q, ctx, testDeviceID, testUserID, "router", "living room")
			seedRecord(t, q, ctx, testRecordID, testDeviceID, "A", "192.0.2.1")

			got, err := q.GetRecordByID(ctx, mustUUID(t, tt.recordID))
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("GetRecordByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetRecordByID() error = %v", err)
			}
			if diff := cmp.Diff(tt.want, got, cmpOptions); diff != "" {
				t.Errorf("GetRecordByID() mismatch (-want +got):\n%s", diff)
			}
			if !got.ReportedAt.Valid {
				t.Errorf("GetRecordByID() reported_at is invalid: %+v", got)
			}
		})
	}
}

func TestListRecordsByDeviceID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		deviceID string
		seed     bool
		wantLen  int
	}{
		{
			name:     "empty",
			deviceID: testDeviceID,
			seed:     false,
			wantLen:  0,
		},
		{
			name:     "success",
			deviceID: testDeviceID,
			seed:     true,
			wantLen:  2,
		},
		{
			name:     "filtered by other device",
			deviceID: testOtherDeviceID,
			seed:     true,
			wantLen:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			if tt.seed {
				seedUser(t, q, ctx, testUserID, "alice", "Alice")
				seedDevice(t, q, ctx, testDeviceID, testUserID, "router", "living room")
				seedRecord(t, q, ctx, testRecordID, testDeviceID, "A", "192.0.2.1")
				seedRecord(t, q, ctx, testOtherRecordID, testDeviceID, "AAAA", "2001:db8::1")
			}

			got, err := q.ListRecordsByDeviceID(ctx, mustUUID(t, tt.deviceID))
			if err != nil {
				t.Fatalf("ListRecordsByDeviceID() error = %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("ListRecordsByDeviceID() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestCreateRecord(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		recordType string
		ipAddress  string
	}{
		{
			name:       "success",
			recordType: "A",
			ipAddress:  "192.0.2.1",
		},
		{
			name:       "null ip address",
			recordType: "TXT",
			ipAddress:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")
			seedDevice(t, q, ctx, testDeviceID, testUserID, "router", "living room")

			params := rdb.CreateRecordParams{
				ID:       mustUUID(t, testRecordID),
				DeviceID: mustUUID(t, testDeviceID),
				Type:     tt.recordType,
			}
			if tt.ipAddress != "" {
				params.IpAddress = mustAddr(t, tt.ipAddress)
			}

			got, err := q.CreateRecord(ctx, params)
			if err != nil {
				t.Fatalf("CreateRecord() error = %v", err)
			}

			want := rdb.Record{
				ID:        mustUUID(t, testRecordID),
				DeviceID:  mustUUID(t, testDeviceID),
				Type:      tt.recordType,
				IpAddress: params.IpAddress,
			}
			if diff := cmp.Diff(want, got, cmpOptions); diff != "" {
				t.Errorf("CreateRecord() mismatch (-want +got):\n%s", diff)
			}
			if !got.ReportedAt.Valid {
				t.Errorf("CreateRecord() reported_at is invalid: %+v", got)
			}

			found, err := q.GetRecordByID(ctx, got.ID)
			if err != nil {
				t.Fatalf("GetRecordByID() after CreateRecord() error = %v", err)
			}
			if diff := cmp.Diff(want, found, cmpOptions); diff != "" {
				t.Errorf("GetRecordByID() after CreateRecord() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateRecord(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		recordID   string
		recordType string
		ipAddress  string
		wantErr    error
	}{
		{
			name:       "success",
			recordID:   testRecordID,
			recordType: "AAAA",
			ipAddress:  "2001:db8::1",
		},
		{
			name:       "null ip address",
			recordID:   testRecordID,
			recordType: "TXT",
		},
		{
			name:       "not found",
			recordID:   testOtherRecordID,
			recordType: "A",
			ipAddress:  "192.0.2.9",
			wantErr:    pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")
			seedDevice(t, q, ctx, testDeviceID, testUserID, "router", "living room")
			created := seedRecord(t, q, ctx, testRecordID, testDeviceID, "A", "192.0.2.1")

			params := rdb.UpdateRecordParams{
				ID:   mustUUID(t, tt.recordID),
				Type: tt.recordType,
			}
			if tt.ipAddress != "" {
				params.IpAddress = mustAddr(t, tt.ipAddress)
			}

			got, err := q.UpdateRecord(ctx, params)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("UpdateRecord() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("UpdateRecord() error = %v", err)
			}

			want := rdb.Record{
				ID:        mustUUID(t, tt.recordID),
				DeviceID:  mustUUID(t, testDeviceID),
				Type:      tt.recordType,
				IpAddress: params.IpAddress,
			}
			if diff := cmp.Diff(want, got, cmpOptions); diff != "" {
				t.Errorf("UpdateRecord() mismatch (-want +got):\n%s", diff)
			}
			if got.ReportedAt.Time != created.ReportedAt.Time {
				t.Errorf("UpdateRecord() changed reported_at: got %v, want %v", got.ReportedAt.Time, created.ReportedAt.Time)
			}
		})
	}
}

func TestDeleteRecord(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		recordID string
	}{
		{
			name:     "success",
			recordID: testRecordID,
		},
		{
			name:     "not found",
			recordID: testOtherRecordID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")
			seedDevice(t, q, ctx, testDeviceID, testUserID, "router", "living room")
			seedRecord(t, q, ctx, testRecordID, testDeviceID, "A", "192.0.2.1")

			if err := q.DeleteRecord(ctx, mustUUID(t, tt.recordID)); err != nil {
				t.Fatalf("DeleteRecord() error = %v", err)
			}

			_, err := q.GetRecordByID(ctx, mustUUID(t, tt.recordID))
			if !errors.Is(err, pgx.ErrNoRows) {
				t.Fatalf("GetRecordByID() after DeleteRecord() error = %v, want %v", err, pgx.ErrNoRows)
			}
		})
	}
}

func TestGetCredentialByUserIDAndProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		userID   string
		provider string
		want     rdb.Credential
		wantErr  error
	}{
		{
			name:     "success",
			userID:   testUserID,
			provider: "local",
			want: rdb.Credential{
				UserID:     mustUUID(t, testUserID),
				Provider:   "local",
				SecretHash: pgtype.Text{String: "hash", Valid: true},
				Subject:    "alice",
			},
		},
		{
			name:     "wrong provider",
			userID:   testUserID,
			provider: "github",
			wantErr:  pgx.ErrNoRows,
		},
		{
			name:     "wrong user",
			userID:   testOtherUserID,
			provider: "local",
			wantErr:  pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")
			seedCredential(t, q, ctx, testUserID, "local", "hash", "alice")

			got, err := q.GetCredentialByUserIDAndProvider(ctx, rdb.GetCredentialByUserIDAndProviderParams{
				UserID:   mustUUID(t, tt.userID),
				Provider: tt.provider,
			})
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("GetCredentialByUserIDAndProvider() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetCredentialByUserIDAndProvider() error = %v", err)
			}
			if diff := cmp.Diff(tt.want, got, cmpOptions); diff != "" {
				t.Errorf("GetCredentialByUserIDAndProvider() mismatch (-want +got):\n%s", diff)
			}
			if !got.CreatedAt.Valid || !got.UpdatedAt.Valid {
				t.Errorf("GetCredentialByUserIDAndProvider() timestamps are invalid: %+v", got)
			}
		})
	}
}

func TestListCredentialsByUserID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  string
		seed    bool
		wantLen int
	}{
		{
			name:    "empty",
			userID:  testUserID,
			seed:    false,
			wantLen: 0,
		},
		{
			name:    "success",
			userID:  testUserID,
			seed:    true,
			wantLen: 2,
		},
		{
			name:    "filtered by other user",
			userID:  testOtherUserID,
			seed:    true,
			wantLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			if tt.seed {
				seedUser(t, q, ctx, testUserID, "alice", "Alice")
				seedCredential(t, q, ctx, testUserID, "local", "hash", "alice")
				seedCredential(t, q, ctx, testUserID, "github", "", "alice-gh")
			}

			got, err := q.ListCredentialsByUserID(ctx, mustUUID(t, tt.userID))
			if err != nil {
				t.Fatalf("ListCredentialsByUserID() error = %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("ListCredentialsByUserID() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestCreateCredential(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		provider   string
		secretHash string
		subject    string
	}{
		{
			name:       "local provider",
			provider:   "local",
			secretHash: "hash",
			subject:    "alice",
		},
		{
			name:     "external provider without secret",
			provider: "github",
			subject:  "alice-gh",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")

			got, err := q.CreateCredential(ctx, rdb.CreateCredentialParams{
				UserID:     mustUUID(t, testUserID),
				Provider:   tt.provider,
				SecretHash: pgtype.Text{String: tt.secretHash, Valid: tt.secretHash != ""},
				Subject:    tt.subject,
			})
			if err != nil {
				t.Fatalf("CreateCredential() error = %v", err)
			}

			want := rdb.Credential{
				UserID:     mustUUID(t, testUserID),
				Provider:   tt.provider,
				SecretHash: pgtype.Text{String: tt.secretHash, Valid: tt.secretHash != ""},
				Subject:    tt.subject,
			}
			if diff := cmp.Diff(want, got, cmpOptions); diff != "" {
				t.Errorf("CreateCredential() mismatch (-want +got):\n%s", diff)
			}
			if !got.CreatedAt.Valid || !got.UpdatedAt.Valid {
				t.Errorf("CreateCredential() timestamps are invalid: %+v", got)
			}

			found, err := q.GetCredentialByUserIDAndProvider(ctx, rdb.GetCredentialByUserIDAndProviderParams{
				UserID:   got.UserID,
				Provider: got.Provider,
			})
			if err != nil {
				t.Fatalf("GetCredentialByUserIDAndProvider() after CreateCredential() error = %v", err)
			}
			if diff := cmp.Diff(want, found, cmpOptions); diff != "" {
				t.Errorf("GetCredentialByUserIDAndProvider() after CreateCredential() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateCredential(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		provider   string
		secretHash string
		subject    string
		wantErr    error
	}{
		{
			name:       "success",
			provider:   "local",
			secretHash: "hash-updated",
			subject:    "alice-updated",
		},
		{
			name:     "not found",
			provider: "github",
			subject:  "nobody",
			wantErr:  pgx.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")
			created := seedCredential(t, q, ctx, testUserID, "local", "hash", "alice")

			got, err := q.UpdateCredential(ctx, rdb.UpdateCredentialParams{
				UserID:     mustUUID(t, testUserID),
				Provider:   tt.provider,
				SecretHash: pgtype.Text{String: tt.secretHash, Valid: tt.secretHash != ""},
				Subject:    tt.subject,
			})
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("UpdateCredential() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("UpdateCredential() error = %v", err)
			}

			want := rdb.Credential{
				UserID:     mustUUID(t, testUserID),
				Provider:   tt.provider,
				SecretHash: pgtype.Text{String: tt.secretHash, Valid: tt.secretHash != ""},
				Subject:    tt.subject,
			}
			if diff := cmp.Diff(want, got, cmpOptions); diff != "" {
				t.Errorf("UpdateCredential() mismatch (-want +got):\n%s", diff)
			}
			if got.CreatedAt.Time != created.CreatedAt.Time {
				t.Errorf("UpdateCredential() changed created_at: got %v, want %v", got.CreatedAt.Time, created.CreatedAt.Time)
			}
			if !got.UpdatedAt.Valid || got.UpdatedAt.Time.Before(created.UpdatedAt.Time) {
				t.Errorf("UpdateCredential() updated_at = %v, want not before %v", got.UpdatedAt.Time, created.UpdatedAt.Time)
			}
		})
	}
}

func TestDeleteCredential(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		userID   string
		provider string
	}{
		{
			name:     "success",
			userID:   testUserID,
			provider: "local",
		},
		{
			name:     "not found",
			userID:   testUserID,
			provider: "github",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q, ctx := newTestQueries(t)
			seedUser(t, q, ctx, testUserID, "alice", "Alice")
			seedCredential(t, q, ctx, testUserID, "local", "hash", "alice")

			err := q.DeleteCredential(ctx, rdb.DeleteCredentialParams{
				UserID:   mustUUID(t, tt.userID),
				Provider: tt.provider,
			})
			if err != nil {
				t.Fatalf("DeleteCredential() error = %v", err)
			}

			_, err = q.GetCredentialByUserIDAndProvider(ctx, rdb.GetCredentialByUserIDAndProviderParams{
				UserID:   mustUUID(t, tt.userID),
				Provider: tt.provider,
			})
			if !errors.Is(err, pgx.ErrNoRows) {
				t.Fatalf("GetCredentialByUserIDAndProvider() after DeleteCredential() error = %v, want %v", err, pgx.ErrNoRows)
			}
		})
	}
}
