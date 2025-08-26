package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)



func TestCheckPasswordHash(t *testing.T) {
	pass1 := "yaypassword!"
	pass2 := "anevenbetterpassword!123"
	hash1, _ := HashPassword(pass1)
	hash2, _ := HashPassword(pass2)

	tests := []struct {
		name 	 string
		password string
		hash	 string
		wantErr  bool
	}{
		{
			name: "Correct password!",
			password: pass1,
			hash: hash1,
			wantErr: false,
		},
		{
			name: "Incorrect password :(",
			password: "bad",
			hash: hash1,
			wantErr: true,
		},
		{
			name: "Mismatched hash",
			password: pass1,
			hash: hash2,
			wantErr: true,
		},
		{
			name: "Empty password",
			password: "",
			hash: hash1,
			wantErr: true,
		},
		{
			name: "Invalid hash",
			password: pass1,
			hash: "shit",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckPasswordHash(test.password, test.hash)
			if (err != nil) != test.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
	
}

func TestValidateJWT(t *testing.T) {
	uid := uuid.New()
	validToken, _ := MakeJWT(uid, "secret", time.Hour)

	tests := []struct {
		name 	 	string
		tokenString string
		tokenSecret	string
		wantUID		uuid.UUID
		wantErr  	bool
	} {
		{
			name: "Valid token!",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUID: uid,
			wantErr: false,
		},
		{
			name: "Invalid token :/",
			tokenString: "shitty",
			tokenSecret: "secret",
			wantUID: uuid.Nil,
			wantErr: true,
		},
		{
			name: "Wrong secret :/",
			tokenString: validToken,
			tokenSecret: "ahhhhhhhh",
			wantUID: uuid.Nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotUID, err := ValidateJWT(test.tokenString, test.tokenSecret)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, test.wantErr)
			}
			if gotUID != test.wantUID {
				t.Errorf("ValidateJWT() gotUID = %v, wantUID %v", gotUID, test.wantUID)
			}
		})
	}

}

func TestGetBearerToken(t *testing.T) {
	head := make(http.Header)
	
	tests := []struct {
		name 	 	string
		setValue	string
		wantToken 	string
		wantErr  	bool
	} {
		{
			name: "Success",
			setValue: "Bearer success",
			wantToken: "success",
			wantErr: false,
		},
		{
			name: "Empty value",
			setValue: "",
			wantToken: "",
			wantErr: true,
		},
	}




	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			head.Set("Authorization", test.setValue)
			gotToken, err := GetBearerToken(head)
			if (err != nil) != test.wantErr {
				t.Errorf("V() error = %v, wantErr %v", err, test.wantErr)
			}
			if gotToken != test.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, wantToken %v", gotToken, test.wantToken)
			}
		})
	}
}