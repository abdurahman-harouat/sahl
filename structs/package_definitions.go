package structs

type PackageDefinition struct {
    PkgName                string   `yaml:"pkgname"`
    PkgVersion             string   `yaml:"pkgver"`
    Dependencies           []string `yaml:"dependencies,omitempty"`
    Source struct {
        URL    string `yaml:"url"`
        MD5    string `yaml:"md5,omitempty"`
        SHA512 string `yaml:"sha512,omitempty"`
    } `yaml:"source"`
    AdditionalDownloads []struct {
        URL    string `yaml:"url"`
        MD5    string `yaml:"md5,omitempty"`
        SHA512 string `yaml:"sha512,omitempty"`
    } `yaml:"additional_downloads,omitempty"`
    Build                       []string `yaml:"build"`
    ExtractedDir                string   `yaml:"extracted_dir,omitempty"`
    Install                     []string `yaml:"install"`
    AdditionalCommands          []string `yaml:"additional_commands,omitempty"`
    AdditionalCommandsWithSudo  []string `yaml:"additional_commands_with_sudo,omitempty"`
}
