package output

// type CardsInBatchesConfig struct {
// 	CardTemplate                    *counters.CardsTemplate
// 	OutputCardTemplate              *counters.CardsTemplate
// 	OutputPathInTemplate            string
// 	GeneratedFileNameStringTemplate string
// }

// CardsToJSONSheets Write files in batches of 70, which is the maximum allowed by TTS (10*7)
// func CardsToJSONSheets(cfg *CardsInBatchesConfig) error {
// 	log.WithFields(log.Fields{"output_file_template": cfg.GeneratedFileNameStringTemplate}).Info("Generating cards files")
//
// 	totalFilesToWrite := int(math.Ceil(float64(len(cfg.OutputCardTemplate.Cards)) / 70))
// 	for fileIteration := 0; fileIteration < totalFilesToWrite; fileIteration++ {
// 		startCard, lastCard := calculateStartAndLastCardIndex(fileIteration, len(cfg.OutputCardTemplate.Cards))
//
// 		cfg.CardTemplate.Cards = nil
// 		cfg.CardTemplate.Cards = cfg.OutputCardTemplate.Cards[startCard:lastCard]
// 		cfg.CardTemplate.OutputPath = fmt.Sprintf(cfg.OutputPathInTemplate, fileIteration+1)
//
// 		generatedOutputFilepath := fmt.Sprintf(cfg.GeneratedFileNameStringTemplate, fileIteration+1)
//
// 		// Write template file
// 		genFile, err := os.OpenFile(generatedOutputFilepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
// 		if err != nil {
// 			log.WithField("filepath", generatedOutputFilepath).Error("error opening file")
// 			return err
// 		}
// 		defer genFile.Close()
//
// 		if err = genFile.Truncate(0); err != nil {
// 			log.WithField("filepath", generatedOutputFilepath).Error("error truncating file")
// 			return err
// 		}
//
// 		byt, _ := json.MarshalIndent(&cfg.CardTemplate, "", "  ")
// 		if _, err = genFile.Write(byt); err != nil {
// 			log.WithField("filepath", generatedOutputFilepath).Error("error writing file")
// 			return err
// 		}
// 	}
//
// 	return nil
// }

// func calculateStartAndLastCardIndex(i, totalCards int) (int, int) {
// 	lastCard := 70 * (i + 1)
// 	if lastCard > totalCards {
// 		lastCard = totalCards
// 	}
// 	startCard := 70 * i
//
// 	return startCard, lastCard
// }
