# Sanitize

Dead simple sanitization.

The goal for this repo is that I often need to be
declarative with my sanitization and this is the best way I've found to do it
for more complex use cases. This has helped me get to the root of my code without
worrying about writing boilerplate validation and test cases

Everything follows this pattern so you can reuse the validator over and over:

```go
validatorStruct.Validate("name of field", value)
validatorStruct.Validate("different field, same rules", otherValue)
```

## Add as dep

```shell
go get github.com/AnthonyHewins/sanitize@latest
```

## Types

Validators

- `NumberValidator[constraints.Integer | constraints.Float]` for integers or floats
- `StrValidator` for strings

Error types (for honing in on errors if you need it)

- `ErrOutOfRange[constraints.Float | constraints.Integer]` for when something falls out of range. This means a string's length falls out of range (in which case it's of type `ErrOutOfRange[int]` because `len(anyString)` is `int`), or a number is out of range: `ErrOutOfRange[float]` for example
- `ErrRequired` when something's required
- `ErrFailedRegexp` when a regexp failed
- `ErrNotUTF8` when something's not UTF8 and needs to be

## Examples

Simple string example

```go
x := sanitize.StrValidator{
    Required: true,
    MinLen: 10,
    MaxLen: 10,
    AllowNonUTF8: false,
    Regex: regexp.MustCompile("^[0-9]{10}$")
}

x.Validate("user phone", "") // &ErrRequired{FieldName: "user phone"}
x.Validate("user phone", "a") // &ErrOutOfRange{FieldName:"user phone", Min: 10, Max: 10, Value: 1}
x.Validate("user phone", "a") // &ErrOutOfRange{FieldName:"user phone", Min: 10, Max: 10, Value: 1}
x.Validate("user phone", "aaaaaaaaaa") // &ErrFailedRegexp{FieldName:"user phone", Regexp: "^[0-9]{10}$"}
x.Validate("user phone", "퟿") // &ErrInvalidUTF8{FieldName: "user phone"}
x.Validate("user phone", "1111111111") // nil
```

Simple int example

```go
hundred := 100
x := sanitize.NumberValidator[int]{
    Required: true,
    MinVal: new(int),
    MaxVal: &hundred,
}

x.ValidatePtr("user's age", nil) // &ErrRequired{FieldName:"user's age"}
x.ValidatePtr("user's age", new(int)) // nil

x.Validate("user's age", 110)// &ErrOutOfRange{FieldName:"user's age", Min: 0, Max: 100, Value: 110}
x.Validate("user's age", -1)// &ErrOutOfRange{FieldName:"user's age", Min: 0, Max: 10, Value: -1}
x.Validate("user's age", 0) // nil
```

More complex example, when you need to do the same validation more than once on a  larger struct. This is the main use case

```go
type Email struct {
    To string
    From string
    CC string
    Body string
}

type EmailValidator struct {
    emailAddress sanitize.StrValidator
    body string
}

func (e EmailValidator) validate(e *Email) error {
    for _,v := range [...][2]string{
        {"to address", e.To}, // re-use this validator over and over
        {"from address", e.From},
        {"CC address", e.CC},
    } {
        if err := e.emailAddress.Validate(v[0], v[1]); err != nil {
            return err
        }
    }

    return e.body.Validate("email body", e.body)
}
```

Now reuse that validator and it will check all emails in your API.
Easier when it's declarative, now your tests are easier:

```go
emailValidatorForAPI := EmailValidator{
    email: sanitize.StrValidator{
        Required: true,
        MinLen: 4,
        MaxLen: 10, // pick to your hearts content
        Regexp: regexp.MustCompile("^[a-z0-9]@[a-z0-9].com$") // not recommended
    },
    body: sanitize.StrValidator{MaxLen: 10000, AllowNonUTF8: true},
}

type AdminBulkEmail struct {
    Emails []Email
    SendWebhooks []*url.URL
}

type UserEmail struct {
    Email Email
}

// ...

userEmail := getUserEmailRequest()
if err := emailValidatorForAPI.validate(&userEmail.Email); err != nil {
    return err
}

adminEmail := getAdminEmailRequest()
for i:=range adminEmail.Emails {
    if err := emailValidatorForAPI.validate(&adminEmail.Emails[i]); err != nil {
        return err
    }
}
```