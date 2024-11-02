# oci-free-arm-instance

Oracle Cloud Infrastructure(OCI) 무료 티어 인스턴스를 생성하는 Go 애플리케이션입니다.

## 요구 사항

- Go 1.22.8 이상
- OCI 계정
- OCI 설정 파일(`~/.oci/config`)
- `.env` 파일

## 설치

1. 저장소를 클론합니다:

    ```sh
    git clone https://github.com/limJiAn/oci-free-arm-instance.git
    cd oci-free-arm-instance
    ```

2. 필요한 Go 모듈을 설치합니다:

    ```sh
    go mod tidy
    ```

3. `.env` 파일을 생성하고 다음 환경 변수를 설정합니다:

    ```env
    OCI_COMPARTMENT_ID=your_compartment_id
    OCI_SUBNET_ID=your_subnet_id
    OCI_IMAGE_ID=your_image_id
    OCI_AVAILABILITY_DOMAIN=your_availability_domain
    OCI_SHAPE=your_shape
    OCI_DISPLAY_NAME=your_display_name

    OCI_BOOT_VOLUME_SIZE_IN_GBS=your_boot_volume_size_in_gbs
    OCI_OCPUS=your_oci_ocpus
    OCI_MEMORY_IN_GBS=your_oci_memory_in_gbs
    ```

4. OCI 설정 파일(`~/.oci/config`)을 생성하고 다음과 같이 설정합니다:

    ```ini
    [DEFAULT]
    user=ocid1.user.oc1.****************
    fingerprint=b1:7a:**:**:**:**:**:**:**:**:**:**:**
    key_file=/home/ubuntu/.oci/oci_api_key.pem
    tenancy=ocid1.tenancy.oc1..****************
    region=*****
    ```

## 사용법

1. 애플리케이션을 실행합니다:

    ```sh
    go run main.go
    ```

2. 성공적으로 실행되면, 콘솔에 생성된 인스턴스의 ID가 출력됩니다.