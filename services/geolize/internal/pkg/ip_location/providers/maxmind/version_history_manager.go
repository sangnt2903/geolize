package maxmind

import (
	"encoding/json"
	"fmt"
	"geolize/services/geolize/internal/pkg/ip_location/model"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type versionHistoryManager struct{}

func newVersionHistoryManager() *versionHistoryManager {
	return &versionHistoryManager{}
}

func (m *versionHistoryManager) CreateFile(payload *model.IPUpdateRequest) (string, error) {
	file := fmt.Sprintf("history__%d__%s.json", time.Now().Unix(), payload.IP)

	var historyFile History = History{
		ID:        payload.IP,
		Name:      payload.IP,
		Overrides: make([]*model.IPUpdateRequest, 0),
	}

	overrideBytes, err := json.Marshal([]*model.IPUpdateRequest{payload})
	if err != nil {
		return file, err
	}

	err = json.Unmarshal(overrideBytes, &historyFile.Overrides)
	if err != nil {
		return file, err
	}

	data, err := json.Marshal(historyFile)
	if err != nil {
		return file, err
	}

	err = os.WriteFile(filepath.Join(dbHistories, file), data, 0644)
	if err != nil {
		return file, err
	}

	return file, nil
}

func (m *versionHistoryManager) RemoveFile(file string) error {
	return os.Remove(filepath.Join(dbHistories, file))
}

func (m *versionHistoryManager) GetAllFiles() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dbHistories, "history*.json"))
	if err != nil {
		return nil, err
	}

	sort.Strings(files)

	return files, nil
}

func (m *versionHistoryManager) GetUpdateFilesFrom(version string) ([]string, error) {
	allFiles, err := m.GetAllFiles()
	if err != nil {
		return nil, err
	}

	if len(allFiles) == 0 {
		return nil, nil
	}

	if len(version) == 0 {
		return allFiles, nil
	}

	_, err = os.Stat(filepath.Join(dbHistories, version))
	if os.IsNotExist(err) {
		return allFiles, nil
	}

	checkpointIndex := -1
	for i, file := range allFiles {
		if filepath.Base(file) == version {
			checkpointIndex = i
			break
		}
	}

	if checkpointIndex == -1 {
		return allFiles, nil
	}

	return allFiles[checkpointIndex:], nil
}

func (m *versionHistoryManager) GetUpdateFilesFromVersion() ([]string, error) {
	version, err := m.GetVersion()
	if err != nil {
		return nil, err
	}

	return m.GetUpdateFilesFrom(version)
}

func (m *versionHistoryManager) GetVersion() (string, error) {
	versionBytes, err := os.ReadFile(versionFilePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(versionBytes)), nil
}

func (m *versionHistoryManager) SetVersion(version string) error {
	return os.WriteFile(versionFilePath, []byte(strings.TrimSpace(version)), 0644)
}
