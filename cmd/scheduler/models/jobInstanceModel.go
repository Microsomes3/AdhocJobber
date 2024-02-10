package models

type JobInstanceModel struct {
	ID                    uint `gorm:"primaryKey"`
	ServerCloudProviderID string
	InstanceStatus        string
	CloudProvider         string
	SSHPublic             []byte
	SSHPrivate            []byte
	RootPassword          string
	IPV4Address           string
	Created               int64         `gorm:"autoCreateTime"` // Use unix seconds as creating time
	JobInstance           TaskRunsModel `gorm:"foreignKey:JobInstanceModelId"`
}
