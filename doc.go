// Package sanitize does basic sanitization of certain data types that can be used to
// compose sanitizers for larger structs.
//
// Int example:
//
//	hundred := 100
//	x := sanitize.NumberValidator[int]{
//	    Required: true,
//	    MinVal: new(int),
//	    MaxVal: &hundred,
//	}
//
//	x.ValidatePtr("user's age", nil) // &ErrRequired{FieldName:"user's age"}
//	x.ValidatePtr("user's age", new(int)) // nil
//
//	x.Validate("user's age", 110)// &ErrOutOfRange{FieldName:"user's age", Min: 0, Max: 100, Value: 110}
//	x.Validate("user's age", -1)// &ErrOutOfRange{FieldName:"user's age", Min: 0, Max: 10, Value: -1}
//	x.Validate("user's age", 0) // nil
//
// String example:
//
//	x := sanitize.StrValidator{
//	    Required: true,
//	    MinLen: 10,
//	    MaxLen: 10,
//	    AllowNonUTF8: false,
//	    Regex: regexp.MustCompile("^[0-9]{10}$")
//	}
//
//	x.Validate("user phone", "") // &ErrRequired{FieldName: "user phone"}
//	x.Validate("user phone", "a") // &ErrOutOfRange{FieldName:"user phone", Min: 10, Max: 10, Value: 1}
//	x.Validate("user phone", "a") // &ErrOutOfRange{FieldName:"user phone", Min: 10, Max: 10, Value: 1}
//	x.Validate("user phone", "aaaaaaaaaa") // &ErrFailedRegexp{FieldName:"user phone", Regexp: "^[0-9]{10}$"}
//	x.Validate("user phone", "퟿") // &ErrInvalidUTF8{FieldName: "user phone"}
//	x.Validate("user phone", "1111111111") // nil
//
// By itself it's not that useful until you start composing the structs yourself. If you have multiple
// checks that could be needed for user input, that's when it'll save you lots of time because you won't
// need to write as much tests:
//
//	type Email struct {
//	    To string
//	    From string
//	    CC string
//	    Body string
//	}
//
//	type EmailValidator struct {
//	    emailAddress sanitize.StrValidator
//	    body string
//	}
//
//	func (e EmailValidator) validate(e *Email) error {
//	    for _,v := range [...][2]string{
//	        {"to address", e.To}, // re-use this validator over and over
//	        {"from address", e.From},
//	        {"CC address", e.CC},
//	    } {
//	        if err := e.emailAddress.Validate(v[0], v[1]); err != nil {
//	            return err
//	        }
//	    }
//
//	    return e.body.Validate("email body", e.body)
//	}
//
// Then use it like so
//
//	emailValidatorForAPI := EmailValidator{
//	    email: sanitize.StrValidator{
//	        Required: true,
//	        MinLen: 4,
//	        MaxLen: 10, // pick to your hearts content
//	        Regexp: regexp.MustCompile("^[a-z0-9]@[a-z0-9].com$") // not recommended
//	    },
//	    body: sanitize.StrValidator{MaxLen: 10000, AllowNonUTF8: true},
//	}
//
//	type AdminBulkEmail struct {
//	    Emails []Email
//	    SendWebhooks []*url.URL
//	}
//
//	type UserEmail struct {
//	    Email Email
//	}
//
//	// ...
//
//	userEmail := getUserEmailRequest()
//	if err := emailValidatorForAPI.validate(&userEmail.Email); err != nil {
//	    return err
//	}
//
//	adminEmail := getAdminEmailRequest()
//	for i:=range adminEmail.Emails {
//	    if err := emailValidatorForAPI.validate(&adminEmail.Emails[i]); err != nil {
//	        return err
//	    }
//	}
package sanitize
