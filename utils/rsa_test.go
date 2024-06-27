package utils

import (
	"testing"
)

func TestRsa(t *testing.T) {
	//data := "LfHbR1fkjdwuvKJNdnHK-Iz1nsH4GY7LxVxuGcceDjSIpikbDPRHjdfD_62-xPNK1PFjiUND5fzCFKo50J0UaHvAGFpF2LOaO3OO8Or4Xp17tlWYBAtRF8hJUeKg6AceeZQXNRw0Rgnb3HIc5_DVcY8jZTZGnXJ3rLuRHRGCqJk"
	rsa := NewRsa("MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDlwS6f4FBSHKDgg8Tti2YXW6ic8BGLeoKI8IuXEUy0q2cV53DcJ7ON55oXuuDuBRLE6PanT86gcoRTp1IOTKjI7fga3arIaWjYubEBzCLUlTPQx/jjO0/mWarj4yvKzk6Ulo/uXWumR+dx0dYiGtbJQlClgILvYtxNHQB7uXWPjwIDAQAB",
		"MIICeQIBADANBgkqhkiG9w0BAQEFAASCAmMwggJfAgEAAoGBAOXBLp/gUFIcoOCDxO2LZhdbqJzwEYt6gojwi5cRTLSrZxXncNwns43nmhe64O4FEsTo9qdPzqByhFOnUg5MqMjt+BrdqshpaNi5sQHMItSVM9DH+OM7T+ZZquPjK8rOTpSWj+5da6ZH53HR1iIa1slCUKWAgu9i3E0dAHu5dY+PAgMBAAECgYEAk87uYeh7g/fq/8WGAZR2v3w2Q5CmmObd559pDm0QvgKvNQZKMzhPaXGgTrfpUPdulcOSOx06vzotK2wvfAeRZUqmApZqlLOiNkcrafEIBjwBlWh7EKxw9bXauKgdQXr7MPfQg11ipbw52wGXmEElvB5tEuCX5tVD9KHzkluXyKECQQD3fq/WgaDTWlnTVb1QvnyBP+bS+40A9JMst9WDK1qKQ8urKFX4Lnfw7s5953Lbx/euLzM1+e9tnWmcUTMa0Op5AkEA7aZtUg48z9rTN4OposMITmOaO870CZot8DE0RS1MshVsSCL6AbKRFiOLzxoDlFGEFtAvephN5qHPtYhWT4bIRwJBAMbY25Ad4EhPjGIWvh9UnJX/8IXNJBIDbwf7v6k+uOTj6YxfwQrA0w8Z34Aa6BabSG2DcMLKR8srMQIt30CJYAkCQQDI6cDWdG75Evkqn8cUcWpeS1qjYa1zSMO5ov+b1FZY4D+xJNDUCpEadGbIaifIhrnzR4I8VPLXHsmpoV/G0B4VAkEArhaCTjjg5KIyyccBIcyTo8RVCQV1/cEwtdl/b+E4JzFatkMvVLbWVSZJ+b0ZxRqDA4DD6qFaZKl2Ya0vgtiPhg==")
	encryptStr := rsa.RsaEncryptStr("sadasdasd")
	t.Log(encryptStr)
}