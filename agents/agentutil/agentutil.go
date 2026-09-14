package agentutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
	"github.com/rs/zerolog/log"
)

const semLogPackageContext = "agent-util::"

// ResolveFilename resolves fn to an absolute path and checks it points to an
// existing regular file. When required is true an empty fn is an error; when
// false an empty fn is returned unchanged with no error.
func ResolveFilename(pn, fn string, required bool) (string, error) {
	const semLogContext = semLogPackageContext + "resolve-filename"
	var err error

	if fn == "" && required {
		err = fmt.Errorf("error: unspecified file %s param", pn)
		log.Error().Err(err).Msg(semLogContext)
		return fn, err
	}

	if fn != "" {
		fn, _ = fileutil.ResolvePath(fn)
		if !fileutil.FileExists(filepath.Dir(fn)) {
			err = fmt.Errorf("error: file %s specifies an invalid path", fn)
			log.Error().Err(err).Msg(semLogContext)
			return fn, err
		}

		if fi, err := os.Stat(fn); err != nil || fi.IsDir() {
			err = fmt.Errorf("error: file %s is not a valid file", fn)
			log.Error().Err(err).Msg(semLogContext)
			return fn, err
		}
	}

	return fn, nil
}

// ResolveFolder resolves folder to an absolute path and checks it points to an
// existing directory. When required is true an empty folder is an error; when
// false an empty folder is returned unchanged with no error.
func ResolveFolder(pn, folder string, required bool) (string, error) {
	const semLogContext = semLogPackageContext + "resolve-folder"
	var err error

	if folder == "" && required {
		err = fmt.Errorf("error: unspecified folder %s param", pn)
		log.Error().Err(err).Msg(semLogContext)
		return folder, err
	}

	if folder != "" {
		folder, _ = fileutil.ResolvePath(folder)
		if !fileutil.FileExists(filepath.Dir(folder)) {
			err = fmt.Errorf("error: folder %s specifies an invalid path", folder)
			log.Error().Err(err).Msg(semLogContext)
			return folder, err
		}

		if fi, err := os.Stat(folder); err != nil || !fi.IsDir() {
			err = fmt.Errorf("error: folder %s is not a valid folder", folder)
			log.Error().Err(err).Msg(semLogContext)
			return folder, err
		}
	}

	return folder, nil
}
