package enaml

type ReleaseManifest struct {
	Version            string                   `yaml:"version,omitempty"`
	Name               string                   `yaml:"name,omitempty"`
	UncommittedChanges bool                     `yaml:"uncommitted_changes,omitempty"`
	CommitHash         string                   `yaml:"commit_hash,omitempty"`
	License            ReleaseLicense           `yaml:"license,omitempty"`
	Jobs               []ReleaseManifestJob     `yaml:"jobs,omitempty"`
	Packages           []ReleaseManifestPackage `yaml:"packages,omitempty"`
}

type ReleaseManifestPackage struct {
	Name         string   `yaml:"name,omitempty"`
	Version      string   `yaml:"version,omitempty"`
	Fingerprint  string   `yaml:"fingerprint,omitempty"`
	SHA1         string   `yaml:"sha1,omitempty"`
	Dependencies []string `yaml:"dependencies"`
}

type ReleaseManifestJob struct {
	Name        string `yaml:"name,omitempty"`
	Version     string `yaml:"version,omitempty"`
	Fingerprint string `yaml:"fingerprint,omitempty"`
	SHA1        string `yaml:"sha1,omitempty"`
}

type ReleaseLicense struct {
	Version     string `yaml:"version,omitempty"`
	Fingerprint string `yaml:"fingerprint,omitempty"`
	SHA1        string `yaml:"sha1,omitempty"`
}
