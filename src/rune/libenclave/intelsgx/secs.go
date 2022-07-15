package intelsgx // import "github.com/confidential-containers/enclave-cc/src/rune/libenclave/intelsgx"

func getSecsAttributes() (uint32, uint32, uint32, uint32) {
	eax, ebx, ecx, edx := cpuid(cpuidSgxFeature, sgxAttributes)
	return eax, ebx, ecx, edx
}
