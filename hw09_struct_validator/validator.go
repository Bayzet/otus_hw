package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrTagNotValid      = errors.New("Tag not valid")
	ErrIsNotStruct      = errors.New("Interface not struct")
	ErrLenNotInt        = errors.New("Rule len not int")
	ErrMinNotInt        = errors.New("Rule min not int")
	ErrMaxNotInt        = errors.New("Rule max not int")
	ErrRegexpNotValid   = errors.New("Regexp not valid")
	ErrIntSlice         = errors.New("Is not integer slice")
	ErrTypeNotSupported = errors.New("Type not supported")

	ErrorLenValidation      = errors.New("Field does not meet the condition len")
	ErrorMinValidation      = errors.New("Field does not meet the condition min")
	ErrorMaxValidation      = errors.New("Field does not meet the condition max")
	ErrorInIntValidation    = errors.New("Field value is not contained in the set")
	ErrorRegexpValidation   = errors.New("Field does not meet the condition regexp")
	ErrorInStringValidation = errors.New("Field value is not contained in the set")
)

const listenTag = "validate"

var rules = []string{"len", "min", "max", "regexp", "in"}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type Condition struct {
	Name  string
	Value interface{}
}

func (c Condition) isValid() bool {
	for _, r := range rules {
		if r == c.Name {
			return true
		}
	}

	return false
}

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range v {
		fmt.Fprintf(&sb, "[Field: %s, Err: %v]", err.Field, err.Err)
	}

	return sb.String()
}

func Validate(v interface{}) error {
	vr := reflect.ValueOf(v)
	var errs ValidationErrors

	if vr.Kind() != reflect.Struct {
		return ErrIsNotStruct
	}

	t := vr.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		ferr := ValidationError{Field: f.Name}

		tag := f.Tag.Get(listenTag)
		if tag == "" {
			continue
		}

		tagsCondition, err := parseTagCondition(tag)
		if err != nil {
			ferr.Err = err
			errs = append(errs, ferr)

			continue
		}

		fv := vr.Field(i)

		switch f.Type.Kind() { //nolint
		case reflect.String:
			err := stringValidator(fv.String(), tagsCondition)
			if err != nil {
				ferr.Err = err
			}
		case reflect.Int:
			err := intValidator(int(fv.Int()), tagsCondition)
			if err != nil {
				ferr.Err = err
			}
		case reflect.Slice:
			err := sliceValidator(fv, tagsCondition)
			if err != nil {
				ferr.Err = err
			}
		}

		if ferr.Err != nil {
			errs = append(errs, ferr)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func parseTagCondition(tag string) (tags []Condition, err error) {
	splitTag := strings.Split(tag, "|")
	if len(splitTag) > 1 {
		for _, v := range splitTag {
			againSplitTag := strings.Split(v, ":")
			c := Condition{againSplitTag[0], againSplitTag[1]}
			if !c.isValid() {
				return nil, ErrTagNotValid
			}
			tags = append(tags, c)
		}

		return
	}

	againSplitTag := strings.Split(tag, ":")
	c := Condition{againSplitTag[0], againSplitTag[1]}
	if !c.isValid() {
		return nil, ErrTagNotValid
	}

	tags = append(tags, c)

	return
}

func stringValidator(v string, c []Condition) (errs error) {
	for _, rule := range c {
		switch rule.Name {
		case "len":
			s := rule.Value.(string)

			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return ErrLenNotInt
			}

			err = lenValidator(v, int(i))
			if err != nil {
				return err
			}
		case "regexp":
			rexp, err := regexp.Compile(rule.Value.(string))
			if err != nil {
				return ErrRegexpNotValid
			}

			err = regexpValidator(v, *rexp)
			if err != nil {
				return err
			}
		case "in":
			if c := strings.Contains(rule.Value.(string), v); !c {
				return ErrorInStringValidation
			}
		}
	}

	return
}

func lenValidator(v string, r int) error {
	if len(v) != r {
		return ErrorLenValidation
	}

	return nil
}

func regexpValidator(v string, r regexp.Regexp) error {
	if r.Find([]byte(v)) == nil {
		return ErrorRegexpValidation
	}

	return nil
}

func intValidator(v int, c []Condition) (errs error) {
	for _, rule := range c {
		switch rule.Name {
		case "min":
			s := rule.Value.(string)
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return ErrMinNotInt
			}

			err = minValidator(v, int(i))
			if err != nil {
				return err
			}
		case "max":
			s := rule.Value.(string)
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return ErrMaxNotInt
			}

			err = maxValidator(v, int(i))
			if err != nil {
				return err
			}
		case "in":
			r := []int{}
			for _, s := range strings.Split(rule.Value.(string), ",") {
				i, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return ErrIntSlice
				}

				r = append(r, int(i))
			}

			err := inIntValidator(v, r)
			if err != nil {
				return err
			}
		}
	}

	return
}

func inIntValidator(v int, r []int) error {
	for _, i := range r {
		if v == i {
			return nil
		}
	}

	return ErrorInIntValidation
}

func minValidator(v int, r int) error {
	if v < r {
		return ErrorMinValidation
	}

	return nil
}

func maxValidator(v int, r int) error {
	if v > r {
		return ErrorMaxValidation
	}

	return nil
}

func sliceValidator(v reflect.Value, c []Condition) error {
	for i := 0; i < v.Len(); i++ {
		switch v.Index(i).Kind() { //nolint
		case reflect.String:
			err := stringValidator(v.Index(i).String(), c)
			if err != nil {
				return err
			}
		case reflect.Int:
			err := intValidator(int(v.Index(i).Int()), c)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
