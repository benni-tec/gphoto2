package gphoto2

import (
    "path"
    "strings"
)

func (cam *Camera) ListFoldersRec(p string) ([]string, error) {
    if !strings.HasSuffix(p, "/") {
        p = p + "/"
    }

    if folders, err := cam.ListFolders(p); err != nil {
        return nil, err
    } else {
        var result []string

        for _, folder := range folders {
            folder = path.Join(p, folder)
            result = append(result, folder)

            if f, err := cam.ListFoldersRec(folder); err != nil {
                return nil, err
            } else {
                for _, f := range f {
                    result = append(result, f)
                }
            }
        }

        return result, nil
    }
}

func (cam *Camera) ListFilesRecYield(p string, yield func(folder string, file string) error) error {
    if !strings.HasSuffix(p, "/") {
        p = p + "/"
    }

    if folders, err := cam.ListFoldersRec(p); err != nil {
        return err
    } else {
        for _, folder := range append([]string{p}, folders...) {
            if f, err := cam.ListFiles(folder); err != nil {
                return err
            } else {
                for _, f := range f {
                    if err := yield(folder, f); err != nil {
                        return err
                    }
                }
            }
        }
    }
    return nil
}

func (cam *Camera) ListFilesRec(p string) ([]string, error) {
    var files []string
    if err := cam.ListFilesRecYield(p, func(folder, file string) error {
        files = append(files, path.Join(folder, file))
        return nil
    }); err != nil {
        return nil, err
    } else {
        return files, nil
    }
}
