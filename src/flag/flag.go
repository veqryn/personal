package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	MyEnv string `arg:"my-env" envconfig:"MY_ENV" required:"true" valid:"required"`
	// MyDuration            time.Duration `default:"60s" split_words:"true"`
}

func main() {

	//blags := flag.NewFlagSet("Flag App", flag.ContinueOnError)
	//blags.SetOutput(os.Stdout)
	//myvar := blags.Int("myvar", 123, "my description")
	//blags.Set("myvar", "234")
	//blags.Parse(os.Args[1:])
	//fmt.Println("----")
	//fmt.Println(*myvar)
	//fmt.Println("----")
	//blags.Usage()
	//fmt.Println("----")
	//blags.PrintDefaults()
	//fmt.Println("----")

	env := Env{}

	blags, err := createFlagSet(os.Stdout, &env)
	if err != nil {
		panic(err)
	}

	err = blags.Parse(os.Args[1:])
	fmt.Println(env)
	if err != nil {
		panic(err)
	}

	//configErr := envconfig.Process("", &env)
	//if configErr != nil {
	//	panic(configErr)
	//}

}

func createFlagSet(output io.Writer, spec interface{}) (*flag.FlagSet, error) {

	blags := flag.NewFlagSet("Flag App", flag.ContinueOnError)
	blags.SetOutput(output)

	s := reflect.ValueOf(spec)

	if s.Kind() != reflect.Ptr {
		return nil, envconfig.ErrInvalidSpecification
	}
	s = s.Elem()
	if s.Kind() != reflect.Struct {
		return nil, envconfig.ErrInvalidSpecification
	}
	typeOfSpec := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		ftype := typeOfSpec.Field(i)
		if !f.CanSet() || isTrue(ftype.Tag.Get("ignored")) {
			continue
		}

		for f.Kind() == reflect.Ptr {
			if f.IsNil() {
				if f.Type().Elem().Kind() != reflect.Struct {
					// nil pointer to a non-struct: leave it alone
					break
				}
				// nil pointer to struct: create a zero instance
				f.Set(reflect.New(f.Type().Elem()))
			}
			f = f.Elem()
		}

		if arg := ftype.Tag.Get("arg"); arg != "" {
			blags.StringVar(f.Addr().Interface().(*string), arg, "", "usage string")
		}
	}

	return blags, nil
}

//
//
//// GatherInfo gathers information about the specified struct
//func gatherInfo(spec interface{}) ([]varInfo, error) {
//	expr := regexp.MustCompile("([^A-Z]+|[A-Z][^A-Z]+|[A-Z]+)")
//	s := reflect.ValueOf(spec)
//
//	if s.Kind() != reflect.Ptr {
//		return nil, envconfig.ErrInvalidSpecification
//	}
//	s = s.Elem()
//	if s.Kind() != reflect.Struct {
//		return nil, envconfig.ErrInvalidSpecification
//	}
//	typeOfSpec := s.Type()
//
//	// over allocate an info array, we will extend if needed later
//	infos := make([]varInfo, 0, s.NumField())
//	for i := 0; i < s.NumField(); i++ {
//		f := s.Field(i)
//		ftype := typeOfSpec.Field(i)
//		if !f.CanSet() || isTrue(ftype.Tag.Get("ignored")) {
//			continue
//		}
//
//		for f.Kind() == reflect.Ptr {
//			if f.IsNil() {
//				if f.Type().Elem().Kind() != reflect.Struct {
//					// nil pointer to a non-struct: leave it alone
//					break
//				}
//				// nil pointer to struct: create a zero instance
//				f.Set(reflect.New(f.Type().Elem()))
//			}
//			f = f.Elem()
//		}
//
//		// Capture information about the config variable
//		info := varInfo{
//			Name:  ftype.Name,
//			Field: f,
//			Tags:  ftype.Tag,
//			Alt:   strings.ToUpper(ftype.Tag.Get("envconfig")),
//		}
//
//		// Default to the field name as the env var name (will be upcased)
//		info.Key = info.Name
//
//		// Best effort to un-pick camel casing as separate words
//		if isTrue(ftype.Tag.Get("split_words")) {
//			words := expr.FindAllStringSubmatch(ftype.Name, -1)
//			if len(words) > 0 {
//				var name []string
//				for _, words := range words {
//					name = append(name, words[0])
//				}
//
//				info.Key = strings.Join(name, "_")
//			}
//		}
//		if info.Alt != "" {
//			info.Key = info.Alt
//		}
//		if prefix != "" {
//			info.Key = fmt.Sprintf("%s_%s", prefix, info.Key)
//		}
//		info.Key = strings.ToUpper(info.Key)
//		infos = append(infos, info)
//
//		if f.Kind() == reflect.Struct {
//			// honor Decode if present
//			if decoderFrom(f) == nil && setterFrom(f) == nil && textUnmarshaler(f) == nil {
//				innerPrefix := prefix
//				if !ftype.Anonymous {
//					innerPrefix = info.Key
//				}
//
//				embeddedPtr := f.Addr().Interface()
//				embeddedInfos, err := gatherInfo(innerPrefix, embeddedPtr)
//				if err != nil {
//					return nil, err
//				}
//				infos = append(infos[:len(infos)-1], embeddedInfos...)
//
//				continue
//			}
//		}
//	}
//	return infos, nil
//}
//
//// varInfo maintains information about the configuration variable
//type varInfo struct {
//	Name  string
//	Alt   string
//	Key   string
//	Field reflect.Value
//	Tags  reflect.StructTag
//}

func isTrue(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}
