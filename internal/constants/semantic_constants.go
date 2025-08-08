package constants

// ========================================
// DELTA-3 SEMANTIC CONSTANTS
// Semantic Analysis and NLP Configuration Constants
// ========================================

// ========================================
// SEMANTIC FIELD NAMES
// ========================================

// Core Semantic Fields
const (
	SemanticFieldNGram       = "ngram"
	SemanticFieldFrequency   = "frequency"
	SemanticFieldScore       = "score"
	SemanticFieldWeight      = "weight"
	SemanticFieldRelevance   = "relevance"
	SemanticFieldSimilarity  = "similarity"
	SemanticFieldDistance    = "distance"
	SemanticFieldCorrelation = "correlation"
	SemanticFieldAnalysis    = "analysis"
	SemanticFieldSemantic    = "semantic"
)

// Text Processing Fields
const (
	SemanticFieldContent     = "content"
	SemanticFieldText        = "text"
	SemanticFieldToken       = "token"
	SemanticFieldWord        = "word"
	SemanticFieldPhrase      = "phrase"
	SemanticFieldSentence    = "sentence"
	SemanticFieldParagraph   = "paragraph"
	SemanticFieldDocument    = "document"
	SemanticFieldTitle       = "title"
	SemanticFieldDescription = "description"
	SemanticFieldKeywords    = "keywords"
	SemanticFieldTags        = "tags"
)

// Semantic Classification Fields
const (
	SemanticFieldCategory       = "category"
	SemanticFieldTopic          = "topic"
	SemanticFieldTheme          = "theme"
	SemanticFieldSentiment      = "sentiment"
	SemanticFieldEmotion        = "emotion"
	SemanticFieldTone           = "tone"
	SemanticFieldContext        = "context"
	SemanticFieldMeaning        = "meaning"
	SemanticFieldConcept        = "concept"
	SemanticFieldEntity         = "entity"
	SemanticFieldRelation       = "relation"
	SemanticFieldAttribute      = "attribute"
	SemanticFieldFeature        = "feature"
)

// Vector and ML Fields
const (
	SemanticFieldVector         = "vector"
	SemanticFieldEmbedding      = "embedding"
	SemanticFieldCluster        = "cluster"
	SemanticFieldClassification = "classification"
	SemanticFieldPrediction     = "prediction"
	SemanticFieldConfidence     = "confidence"
	SemanticFieldProbability    = "probability"
	SemanticFieldThreshold      = "threshold"
)

// Statistical Fields
const (
	SemanticFieldMinimum   = "minimum"
	SemanticFieldMaximum   = "maximum"
	SemanticFieldAverage   = "average"
	SemanticFieldMedian    = "median"
	SemanticFieldVariance  = "variance"
	SemanticFieldDeviation = "deviation"
)

// ========================================
// SEMANTIC ANALYSIS TYPES
// ========================================

// N-Gram Analysis Types
const (
	SemanticAnalysisUnigram   = "unigram"
	SemanticAnalysisBigram    = "bigram"
	SemanticAnalysisTrigram   = "trigram"
	SemanticAnalysisFourGram  = "4-gram"
	SemanticAnalysisFiveGram  = "5-gram"
	SemanticAnalysisNGram     = "n-gram"
)

// Similarity Metrics
const (
	SemanticAnalysisTFIDF       = "tf-idf"
	SemanticAnalysisCosine      = "cosine"
	SemanticAnalysisEuclidean   = "euclidean"
	SemanticAnalysisJaccard     = "jaccard"
	SemanticAnalysisLevenshtein = "levenshtein"
	SemanticAnalysisHamming     = "hamming"
	SemanticAnalysisJaro        = "jaro"
	SemanticAnalysisSoundex     = "soundex"
	SemanticAnalysisMetaphone   = "metaphone"
)

// Text Processing Types
const (
	SemanticAnalysisStemming       = "stemming"
	SemanticAnalysisLemmatization  = "lemmatization"
	SemanticAnalysisTokenization   = "tokenization"
	SemanticAnalysisNormalization  = "normalization"
	SemanticAnalysisPreprocessing  = "preprocessing"
	SemanticAnalysisPostprocessing = "postprocessing"
	SemanticAnalysisExtraction     = "extraction"
	SemanticAnalysisEnrichment     = "enrichment"
	SemanticAnalysisAnnotation     = "annotation"
	SemanticAnalysisTagging        = "tagging"
	SemanticAnalysisParsing        = "parsing"
	SemanticAnalysisSegmentation   = "segmentation"
)

// ML Analysis Types
const (
	SemanticAnalysisClustering      = "clustering"
	SemanticAnalysisClassification  = "classification"
	SemanticAnalysisRegression      = "regression"
	SemanticAnalysisDimensionality  = "dimensionality"
	SemanticAnalysisReduction       = "reduction"
	SemanticAnalysisOptimization    = "optimization"
)

// ========================================
// NLP CONSTANTS
// ========================================

// Character and Encoding Constants
const (
	SemanticNLPStopwords    = "stopwords"
	SemanticNLPPunctuation  = "punctuation"
	SemanticNLPWhitespace   = "whitespace"
	SemanticNLPAlphanumeric = "alphanumeric"
	SemanticNLPUnicode      = "unicode"
	SemanticNLPASCII        = "ascii"
	SemanticNLPUTF8         = "utf8"
	SemanticNLPEncoding     = "encoding"
)

// Language Processing Constants
const (
	SemanticNLPLanguage    = "language"
	SemanticNLPLocale      = "locale"
	SemanticNLPGrammar     = "grammar"
	SemanticNLPSyntax      = "syntax"
	SemanticNLPMorphology  = "morphology"
	SemanticNLPPhonetics   = "phonetics"
	SemanticNLPSemantics   = "semantics"
	SemanticNLPPragmatics  = "pragmatics"
	SemanticNLPDiscourse   = "discourse"
)

// Knowledge Base Constants
const (
	SemanticNLPCorpus     = "corpus"
	SemanticNLPLexicon    = "lexicon"
	SemanticNLPVocabulary = "vocabulary"
	SemanticNLPDictionary = "dictionary"
	SemanticNLPThesaurus  = "thesaurus"
	SemanticNLPOntology   = "ontology"
	SemanticNLPTaxonomy   = "taxonomy"
	SemanticNLPHierarchy  = "hierarchy"
)

// Data Structure Constants
const (
	SemanticNLPGraph   = "graph"
	SemanticNLPTree    = "tree"
	SemanticNLPNetwork = "network"
	SemanticNLPMatrix  = "matrix"
	SemanticNLPTensor  = "tensor"
	SemanticNLPArray   = "array"
	SemanticNLPList    = "list"
	SemanticNLPSet     = "set"
	SemanticNLPMap     = "map"
	SemanticNLPHash    = "hash"
	SemanticNLPIndex   = "index"
	SemanticNLPCache   = "cache"
	SemanticNLPBuffer  = "buffer"
	SemanticNLPQueue   = "queue"
	SemanticNLPStack   = "stack"
	SemanticNLPHeap    = "heap"
)

// ========================================
// SEMANTIC METRICS
// ========================================

// Classification Metrics
const (
	SemanticMetricPrecision = "precision"
	SemanticMetricRecall    = "recall"
	SemanticMetricF1        = "f1"
	SemanticMetricF2        = "f2"
	SemanticMetricAccuracy  = "accuracy"
	SemanticMetricAUC       = "auc"
	SemanticMetricROC       = "roc"
)

// Regression Metrics
const (
	SemanticMetricMAE         = "mae"
	SemanticMetricMSE         = "mse"
	SemanticMetricRMSE        = "rmse"
	SemanticMetricR2          = "r2"
	SemanticMetricCorrelation = "correlation"
	SemanticMetricCovariance  = "covariance"
)

// Information Theory Metrics
const (
	SemanticMetricEntropy           = "entropy"
	SemanticMetricMutualInformation = "mutual_information"
	SemanticMetricKLDivergence      = "kl_divergence"
	SemanticMetricJSDivergence      = "js_divergence"
)

// Statistical Test Metrics
const (
	SemanticMetricChiSquare = "chi_square"
	SemanticMetricPValue    = "p_value"
	SemanticMetricZScore    = "z_score"
	SemanticMetricTTest     = "t_test"
	SemanticMetricANOVA     = "anova"
)

// ========================================
// SEMANTIC ALGORITHMS
// ========================================

// Traditional ML Algorithms
const (
	SemanticAlgorithmNaiveBayes       = "naive_bayes"
	SemanticAlgorithmSVM              = "svm"
	SemanticAlgorithmKNN              = "knn"
	SemanticAlgorithmDecisionTree     = "decision_tree"
	SemanticAlgorithmRandomForest     = "random_forest"
	SemanticAlgorithmGradientBoosting = "gradient_boosting"
	SemanticAlgorithmXGBoost          = "xgboost"
	SemanticAlgorithmLightGBM         = "lightgbm"
)

// Neural Network Algorithms
const (
	SemanticAlgorithmNeuralNetwork = "neural_network"
	SemanticAlgorithmDeepLearning  = "deep_learning"
	SemanticAlgorithmCNN           = "cnn"
	SemanticAlgorithmRNN           = "rnn"
	SemanticAlgorithmLSTM          = "lstm"
	SemanticAlgorithmGRU           = "gru"
	SemanticAlgorithmTransformer   = "transformer"
)

// Pre-trained Models
const (
	SemanticAlgorithmBERT     = "bert"
	SemanticAlgorithmGPT      = "gpt"
	SemanticAlgorithmWord2Vec = "word2vec"
	SemanticAlgorithmGlove    = "glove"
	SemanticAlgorithmFastText = "fasttext"
	SemanticAlgorithmDoc2Vec  = "doc2vec"
)

// Topic Modeling Algorithms
const (
	SemanticAlgorithmLDA = "lda"
	SemanticAlgorithmLSA = "lsa"
	SemanticAlgorithmNMF = "nmf"
)

// Dimensionality Reduction Algorithms
const (
	SemanticAlgorithmPCA  = "pca"
	SemanticAlgorithmTSNE = "tsne"
	SemanticAlgorithmUMAP = "umap"
)

// Clustering Algorithms
const (
	SemanticAlgorithmDBSCAN           = "dbscan"
	SemanticAlgorithmKMeans           = "kmeans"
	SemanticAlgorithmHierarchical     = "hierarchical"
	SemanticAlgorithmSpectral         = "spectral"
	SemanticAlgorithmGaussianMixture  = "gaussian_mixture"
)

// Anomaly Detection Algorithms
const (
	SemanticAlgorithmIsolationForest    = "isolation_forest"
	SemanticAlgorithmOneClassSVM        = "one_class_svm"
	SemanticAlgorithmLocalOutlierFactor = "local_outlier_factor"
)

// ========================================
// SEMANTIC DATA TYPES
// ========================================

// Basic Data Types
const (
	SemanticDataTypeString     = "string"
	SemanticDataTypeText       = "text"
	SemanticDataTypeNumeric    = "numeric"
	SemanticDataTypeInteger    = "integer"
	SemanticDataTypeFloat      = "float"
	SemanticDataTypeDouble     = "double"
	SemanticDataTypeBoolean    = "boolean"
)

// Categorical Data Types
const (
	SemanticDataTypeCategorical = "categorical"
	SemanticDataTypeOrdinal     = "ordinal"
	SemanticDataTypeNominal     = "nominal"
	SemanticDataTypeBinary      = "binary"
)

// Matrix Types
const (
	SemanticDataTypeSparse = "sparse"
	SemanticDataTypeDense  = "dense"
)

// Structure Types
const (
	SemanticDataTypeStructured   = "structured"
	SemanticDataTypeUnstructured = "unstructured"
)

// Learning Types
const (
	SemanticDataTypeLabeled        = "labeled"
	SemanticDataTypeUnlabeled      = "unlabeled"
	SemanticDataTypeSupervised     = "supervised"
	SemanticDataTypeUnsupervised   = "unsupervised"
	SemanticDataTypeSemiSupervised = "semi_supervised"
	SemanticDataTypeReinforcement  = "reinforcement"
)

// Processing Types
const (
	SemanticDataTypeOnline      = "online"
	SemanticDataTypeOffline     = "offline"
	SemanticDataTypeBatch       = "batch"
	SemanticDataTypeStream      = "stream"
	SemanticDataTypeRealTime    = "real_time"
	SemanticDataTypeNearRealTime = "near_real_time"
)

// ========================================
// SEMANTIC CONFIG KEYS
// ========================================

// Model Configuration
const (
	SemanticConfigModel           = "model"
	SemanticConfigAlgorithm       = "algorithm"
	SemanticConfigParameters      = "parameters"
	SemanticConfigHyperparameters = "hyperparameters"
	SemanticConfigFeatures        = "features"
	SemanticConfigDimensions      = "dimensions"
)

// Neural Network Configuration
const (
	SemanticConfigLayers        = "layers"
	SemanticConfigNeurons       = "neurons"
	SemanticConfigEpochs        = "epochs"
	SemanticConfigBatchSize     = "batch_size"
	SemanticConfigLearningRate  = "learning_rate"
	SemanticConfigDropout       = "dropout"
	SemanticConfigRegularization = "regularization"
	SemanticConfigOptimizer     = "optimizer"
	SemanticConfigLoss          = "loss"
	SemanticConfigMetric        = "metric"
	SemanticConfigThreshold     = "threshold"
	SemanticConfigWindow        = "window"
	SemanticConfigStride        = "stride"
	SemanticConfigPadding       = "padding"
	SemanticConfigActivation    = "activation"
	SemanticConfigInitializer   = "initializer"
	SemanticConfigScheduler     = "scheduler"
)

// Training Configuration
const (
	SemanticConfigEarlyStopping = "early_stopping"
	SemanticConfigCheckpointing = "checkpointing"
	SemanticConfigLogging       = "logging"
	SemanticConfigDebugging     = "debugging"
	SemanticConfigProfiling     = "profiling"
	SemanticConfigMonitoring    = "monitoring"
)

// ========================================
// PREPROCESSING STEPS
// ========================================

// Data Cleaning
const (
	SemanticPreprocessCleaning       = "cleaning"
	SemanticPreprocessNormalization  = "normalization"
	SemanticPreprocessStandardization = "standardization"
	SemanticPreprocessScaling        = "scaling"
)

// Data Transformation
const (
	SemanticPreprocessEncoding        = "encoding"
	SemanticPreprocessDecoding        = "decoding"
	SemanticPreprocessTransformation  = "transformation"
	SemanticPreprocessAugmentation    = "augmentation"
)

// Data Sampling
const (
	SemanticPreprocessSampling  = "sampling"
	SemanticPreprocessFiltering = "filtering"
	SemanticPreprocessSelection = "selection"
	SemanticPreprocessExtraction = "extraction"
	SemanticPreprocessReduction = "reduction"
)

// Data Compression
const (
	SemanticPreprocessCompression   = "compression"
	SemanticPreprocessDecompression = "decompression"
	SemanticPreprocessEncryption    = "encryption"
	SemanticPreprocessDecryption    = "decryption"
)

// Data Organization
const (
	SemanticPreprocessHashing     = "hashing"
	SemanticPreprocessIndexing    = "indexing"
	SemanticPreprocessSorting     = "sorting"
	SemanticPreprocessShuffling   = "shuffling"
	SemanticPreprocessSplitting   = "splitting"
	SemanticPreprocessMerging     = "merging"
	SemanticPreprocessJoining     = "joining"
	SemanticPreprocessAggregation = "aggregation"
	SemanticPreprocessGrouping    = "grouping"
)

// ========================================
// FILE FORMATS
// ========================================

// Text File Formats
const (
	SemanticFileFormatTXT  = ".txt"
	SemanticFileFormatCSV  = ".csv"
	SemanticFileFormatJSON = ".json"
	SemanticFileFormatXML  = ".xml"
	SemanticFileFormatHTML = ".html"
	SemanticFileFormatMD   = ".md"
	SemanticFileFormatYAML = ".yaml"
	SemanticFileFormatYML  = ".yml"
)

// Document File Formats
const (
	SemanticFileFormatPDF  = ".pdf"
	SemanticFileFormatDOC  = ".doc"
	SemanticFileFormatDOCX = ".docx"
	SemanticFileFormatRTF  = ".rtf"
	SemanticFileFormatODT  = ".odt"
)

// ========================================
// DEFAULT VALUES
// ========================================

// Default Analysis Parameters
const (
	SemanticDefaultNGramSize     = 3
	SemanticDefaultMinFrequency  = 2
	SemanticDefaultMaxFeatures   = 10000
	SemanticDefaultThreshold     = 0.5
	SemanticDefaultSimilarity    = 0.8
	SemanticDefaultConfidence    = 0.95
)

// Default Processing Parameters
const (
	SemanticDefaultBatchSize    = 100
	SemanticDefaultTimeout      = 300 // seconds
	SemanticDefaultMaxRetries   = 3
	SemanticDefaultCacheSize    = 1000
	SemanticDefaultBufferSize   = 4096
)

// Default Model Parameters
const (
	SemanticDefaultEpochs       = 100
	SemanticDefaultLearningRate = 0.001
	SemanticDefaultDropout      = 0.2
	SemanticDefaultL2Reg        = 0.01
)