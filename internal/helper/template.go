package helper

import "html/template"

func LoadTemplate(tmplName string) (*template.Template, error) {
    files := []string{
        "web/templates/base.html",
        "web/templates/" + tmplName + ".html",
    }

    t, err := template.ParseFiles(files...)
    if err != nil {
        return nil, err 
    }

    return t, nil
}
