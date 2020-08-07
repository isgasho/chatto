package clf

import (
	"log"
	"strings"

	"github.com/navossoc/bayesian"
	"github.com/spf13/viper"
)

// Classification models a classification yaml file
type Classification struct {
	Classification []TrainingTexts `yaml:"classification"`
}

// TrainingTexts models texts used for training the classifier
type TrainingTexts struct {
	Command string   `yaml:"command"`
	Texts   []string `yaml:"texts"`
}

// Classifier models a classifier and its classes
type Classifier struct {
	Model   bayesian.Classifier
	Classes []bayesian.Class
}

// Predict predict a class for a given text
func (c *Classifier) Predict(text string) (string, float64) {
	probs, likely, _ := c.Model.ProbScores(Clean(&text))
	class := string(c.Classes[likely])
	prob := probs[likely]
	// log.Printf("CLF\t| \"%v\" classified as %v (%0.2f%%)", text, class, prob*100)
	// if prob < 0.75/float64(len(probs)) {
	// 	return "", -1.0
	// }
	return class, prob
}

// Load loads classification configuration from yaml
func Load(path *string) Classification {
	config := viper.New()
	config.SetConfigName("clf")
	config.AddConfigPath(*path)

	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	var botClassif Classification
	if err := config.Unmarshal(&botClassif); err != nil {
		panic(err)
	}

	return botClassif
}

// Create returns a trained Classifier
func Create(path *string) Classifier {
	classification := Load(path)

	// classes := make([]bayesian.Class, 0)
	var classes []bayesian.Class
	for _, class := range classification.Classification {
		classes = append(classes, bayesian.Class(class.Command))
	}

	classifier := bayesian.NewClassifier(classes...)

	for _, cls := range classification.Classification {
		for _, txt := range cls.Texts {
			classifier.Learn(Clean(&txt), bayesian.Class(cls.Command))
		}
	}

	log.Println("Loaded commands for classifier:")
	for i, c := range classes {
		log.Printf("%v\t%v\n", i, c)
	}

	return Classifier{*classifier, classes}
}

// Clean performs steps to clean a string
func Clean(text *string) []string {
	tokens := strings.Split(*text, " ")
	lowerTokens := make([]string, 0)
	for _, t := range tokens {
		lowerTokens = append(lowerTokens, strings.ToLower(t))
	}
	return lowerTokens
}
