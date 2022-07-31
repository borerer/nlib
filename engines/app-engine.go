package engines

type AppEngine struct {
	AppID      string
	FileEngine FileEngine
}

func NewAppEngine(appID string, fileBaseDir string) *AppEngine {
	appEngine := &AppEngine{
		AppID: appID,
	}
	appEngine.FileEngine = &FileFS{
		BaseDir: fileBaseDir,
	}
	return appEngine
}
