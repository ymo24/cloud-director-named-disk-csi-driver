package e2e

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware/cloud-provider-for-cloud-director/pkg/testingsdk"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	apiv1 "k8s.io/api/core/v1"
	"os"
	"testing"
)

func init() {
	rdeID = os.Getenv("RdeID")
	if rdeID == "" {
		// It is okay to panic here as this will be caught during dev
		panic("RdeID should be set")
	}
	vcdHost = os.Getenv("VcdHost")
	if vcdHost == "" {
		// It is okay to panic here as this will be caught during dev
		panic("VcdHost should be set")
	}
	orgName = os.Getenv("OrgName")
	if orgName == "" {
		// It is okay to panic here as this will be caught during dev
		panic("OrgName should be set")
	}
	ovdc = os.Getenv("Ovdc")
	if ovdc == "" {
		// It is okay to panic here as this will be caught during dev
		panic("Ovdc should be set")
	}
	userName = os.Getenv("UserName")
	if userName == "" {
		// It is okay to panic here as this will be caught during dev
		panic("UserName should be set")
	}
	userOrg = os.Getenv("UserOrg")
	if userOrg == "" {
		// It is okay to panic here as this will be caught during dev
		panic("UserOrg should be set")
	}
	refreshToken = os.Getenv("RefreshToken")
	if refreshToken == "" {
		// It is okay to panic here as this will be caught during dev
		panic("RefreshToken should be set")
	}
}

const (
	staticPVExceedLimit     = "disk-size-exceed-limit"
	DiskName                = "csi-test-disk"
	staticPVCName           = "static-pv-claim"
	dynamicPVCName          = "dynamic-pv-claim"
	storageClass            = "default-storage-class-1"
	storageSize             = "2Gi"
	storageSizeMB           = 2048
	deployName              = "nginx-deployment"
	kubeConfigPath          = "/Users/ymo/Downloads/kubeconfig-csi-auto.txt"
	storageProfileWithLimit = "development2"
	defaultStorageProfile   = "*"
)

var (
	rdeID        string
	vcdHost      string
	orgName      string
	ovdc         string
	userName     string
	userOrg      string
	refreshToken string
)

func TestCSIAutomation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CSI Testing Suite")
}

var _ = Describe("CSI automation Test", func() {
	Context("CSI automation test", func() {
		var (
			testingClient *testingsdk.TestClient
			//diskClient    *DiskClient
		)
		const nameSpace4 = "name-space-test-2"
		const staticDisk1 = "static-disk-delete-6"
		const staticDisk2 = "static-disk-retain"
		const DeletestorageClass = "delete-sc"
		const RetainstorageClass = "retain-storage-class"
		const staticPVCNameRetain = "pvc-static-retain"
		const staticPVCNameDelete = "pvc-static-delete"
		const dynamicPVCNameDelete = "pvc-dynamic-delete"
		const dynamicPVCNameRetain = "pvc-dynamic-retain"
		const volumeName = "deployment-pv"

		BeforeEach(func() {
			var err error
			testingClient, err = testingsdk.NewTestClient(&testingsdk.VCDAuthParams{
				Host:         vcdHost,
				OvdcName:     ovdc,
				OrgName:      orgName,
				Username:     userName,
				RefreshToken: refreshToken,
				UserOrg:      userOrg,
				GetVdcClient: true,
			}, rdeID)
			Expect(err).NotTo(HaveOccurred())
			ns, err := testingClient.CreateNameSpace(context.TODO(), nameSpace4)
			Expect(err).NotTo(HaveOccurred())
			Expect(ns).NotTo(BeNil())
			retainStorageClass, err := testingClient.CreateStorageClass(context.TODO(), RetainstorageClass, apiv1.PersistentVolumeReclaimRetain, defaultStorageProfile)
			Expect(err).NotTo(HaveOccurred())
			Expect(retainStorageClass).NotTo(BeNil())
			deleteStorageClass, err := testingClient.CreateStorageClass(context.TODO(), DeletestorageClass, apiv1.PersistentVolumeReclaimDelete, defaultStorageProfile)
			Expect(err).NotTo(HaveOccurred())
			Expect(deleteStorageClass).NotTo(BeNil())
		})
		//It("PV static provisioning test using delete reclaim policy", func() {
		//	err := CreateDisk(testingClient.VcdClient, staticDisk1, storageSizeMB)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("Disk is created successfully from VCD")
		//	pv, err := testingClient.CreatePV(context.TODO(), staticDisk1, storageClassDelete, defaultStorageProfile, storageSize, apiv1.PersistentVolumeReclaimDelete)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pv).NotTo(BeNil())
		//	pv, err = testingClient.GetPV(context.TODO(), staticDisk1)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pv).NotTo(BeNil())
		//	By("PV is created successfully from VCD")
		//	pvc, err := testingClient.CreatePVC(context.TODO(), nameSpace4, staticPVCNameDelete, storageClassDelete, storageSize)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pvc).NotTo(BeNil())
		//	err = testingClient.IsPvcReady(context.TODO(), nameSpace4, staticPVCNameDelete)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("PVC is created successfully")
		//	deployment, err := testingClient.CreateDeployment(context.TODO(), &testingsdk.DeployParams{
		//		Name: deployName,
		//		Labels: map[string]string{
		//			"app": "nginx",
		//		},
		//		ContainerParams: testingsdk.ContainerParams{
		//			ContainerName:  "nginx",
		//			ContainerImage: "nginx:1.14.2",
		//			ContainerPort:  80,
		//		},
		//		VolumeParams: testingsdk.VolumeParams{
		//			VolumeName: staticDisk1,
		//			PvcRef:     staticPVCNameDelete,
		//			MountPath:  "/init-container-msg-mount-path",
		//		},
		//	}, nameSpace4)
		//	Expect(deployment).NotTo(BeNil())
		//	By("Deployment is created successfully")
		//	err = testingClient.IsDeploymentReady(context.TODO(), nameSpace4, dynamicPVCNameDelete)
		//	Expect(err).NotTo(HaveOccurred())
		//
		//	By("Deployment is ready")
		//	err = testingClient.DeletePVC(context.TODO(), nameSpace4, staticPVCNameDelete)
		//	By("PVC is deleted successfully")
		//	err = testingClient.DeleteDeployment(context.TODO(), nameSpace4, deployName)
		//	By("Deployment is deleted successfully")
		//	pv, err = testingClient.GetPV(context.TODO(), staticDisk1)
		//	By("PV is presented in Kubernetes")
		//	vcdDisk, err := getDiskByNameViaVCD(testingClient.VcdClient, staticDisk1)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(vcdDisk).NotTo(BeNil())
		//	By("PV is verified in VCD")
		//	pvFound, err := GetPVByNameViaRDE(staticDisk1, testingClient, "named-disk")
		//	Expect(err).To(MatchError(testingsdk.ResourceNotFound))
		//	Expect(pvFound).To(BeFalse())
		//	By("PV is not shown in RDE")
		//	err = testingClient.DeletePV(context.TODO(), staticDisk1)
		//	Expect(err).NotTo(HaveOccurred())
		//	//Todo: Restore above
		//	vcdDisk, err = getDiskByNameViaVCD(testingClient.VcdClient, staticDisk1)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(vcdDisk).NotTo(BeNil())
		//	pv, err = testingClient.GetPV(context.TODO(), staticDisk1)
		//	Expect(err).To(HaveOccurred())
		//	Expect(pv).To(BeNil())
		//	By("PV is deleted in Kubernetes")
		//	pvFound, err = GetPVByNameViaRDE(staticDisk1, testingClient, "named-disk")
		//	Expect(err).To(MatchError(testingsdk.ResourceNotFound))
		//	Expect(pvFound).To(BeFalse())
		//	By("PV is not shown in RDE")
		//	err = DeleteDisk(testingClient.VcdClient, staticDisk1)
		//	vcdDisk, err = getDiskByNameViaVCD(testingClient.VcdClient, staticDisk1)
		//	Expect(err).To(MatchError(govcd.ErrorEntityNotFound))
		//	Expect(vcdDisk).To(BeNil())
		//
		//})

		It("Testing Verification test", func() {
			adminOrg, err := testingClient.VcdClient.VCDClient.GetAdminOrgByName(orgName)
			Expect(err).NotTo(HaveOccurred())
			Expect(adminOrg).NotTo(BeNil())
			adminVdc, err := adminOrg.GetAdminVDCByName(ovdc, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(adminVdc).NotTo(BeNil())
			rawSpList, err := testingClient.VcdClient.VCDClient.Client.QueryAllProviderVdcStorageProfiles()
			var spRecord *types.QueryResultProviderVdcStorageProfileRecordType
			for _, sp := range rawSpList {
				if sp.Name == storageProfileWithLimit {
					spRecord = sp
				}
			}
			err = adminVdc.AddStorageProfileWait(&types.VdcStorageProfileConfiguration{
				Enabled: true,
				Units:   "MB",
				Limit:   2048,
				Default: false,
				ProviderVdcStorageProfile: &types.Reference{
					HREF: spRecord.HREF,
					Name: spRecord.Name,
				},
			}, "storage profile for ")
			Expect(err).NotTo(HaveOccurred())
			err = CreateDisk(testingClient.VcdClient, staticDisk1, 4096, storageProfileWithLimit)
			Expect(err).To(HaveOccurred())
			isQuotaError := validateDiskQuotaError(err)
			Expect(isQuotaError).To(BeTrue())
			err = adminVdc.RemoveStorageProfileWait(storageProfileWithLimit)
			Expect(err).NotTo(HaveOccurred())
		})
		//It("PV static provisioning test using retain reclaim policy", func() {
		//	err := CreateDisk(testingClient.VcdClient, staticDisk1, storageSizeMB)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("Disk is created successfully from VCD")
		//	pv, err := testingClient.CreatePV(context.TODO(), staticDisk1, storageClassRetain, defaultStorageProfile, storageSize, apiv1.PersistentVolumeReclaimRetain)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pv).NotTo(BeNil())
		//	pv, err = testingClient.GetPV(context.TODO(), staticDisk1)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pv).NotTo(BeNil())
		//	By("PV is created successfully from VCD")
		//	pvc, err := testingClient.CreatePVC(context.TODO(), nameSpace4, staticPVCNameRetain, storageClassRetain, storageSize)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pvc).NotTo(BeNil())
		//	err = testingClient.IsPvcReady(context.TODO(), nameSpace4, staticPVCNameRetain)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("PVC is created successfully")
		//	deployment, err := testingClient.CreateDeployment(context.TODO(), &testingsdk.DeployParams{
		//		Name: deployName,
		//		Labels: map[string]string{
		//			"app": "nginx",
		//		},
		//		ContainerParams: testingsdk.ContainerParams{
		//			ContainerName:  "nginx",
		//			ContainerImage: "nginx:1.14.2",
		//			ContainerPort:  80,
		//		},
		//		VolumeParams: testingsdk.VolumeParams{
		//			VolumeName: staticDisk1,
		//			PvcRef:     staticPVCNameRetain,
		//			MountPath:  "/init-container-msg-mount-path",
		//		},
		//	}, nameSpace4)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(deployment).NotTo(BeNil())
		//	By("Deployment is created successfully")
		//	err = testingClient.IsDeploymentReady(context.TODO(), nameSpace4, dynamicPVCNameRetain)
		//	Expect(err).NotTo(HaveOccurred())
		//
		//	By("Deployment is ready")
		//	err = testingClient.DeletePVC(context.TODO(), nameSpace4, staticPVCNameRetain)
		//	By("PVC is deleted successfully")
		//	err = testingClient.DeleteDeployment(context.TODO(), nameSpace4, deployName)
		//	By("Deployment is deleted successfully")
		//	pv, err = testingClient.GetPV(context.TODO(), staticDisk1)
		//	By("PV is presented in Kubernetes")
		//	vcdDisk, err := getDiskByNameViaVCD(testingClient.VcdClient, staticDisk1)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(vcdDisk).NotTo(BeNil())
		//	By("PV is verified in VCD")
		//	pvFound, err := GetPVByNameViaRDE(staticDisk1, testingClient, "named-disk")
		//	Expect(err).To(MatchError(testingsdk.ResourceNotFound))
		//	Expect(pvFound).To(BeFalse())
		//	By("PV is not shown in RDE")
		//	err = testingClient.DeletePV(context.TODO(), staticDisk1)
		//	Expect(err).NotTo(HaveOccurred())
		//	//Todo: Restore above
		//	vcdDisk, err = getDiskByNameViaVCD(testingClient.VcdClient, staticDisk1)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(vcdDisk).NotTo(BeNil())
		//	pv, err = testingClient.GetPV(context.TODO(), staticDisk1)
		//	Expect(err).To(HaveOccurred())
		//	Expect(pv).To(BeNil())
		//	By("PV is deleted in Kubernetes")
		//	pvFound, err = GetPVByNameViaRDE(staticDisk1, testingClient, "named-disk")
		//	Expect(err).To(MatchError(testingsdk.ResourceNotFound))
		//	Expect(pvFound).To(BeFalse())
		//	By("PV is not shown in RDE")
		//	err = DeleteDisk(testingClient.VcdClient, staticDisk1)
		//	vcdDisk, err = getDiskByNameViaVCD(testingClient.VcdClient, staticDisk1)
		//	Expect(err).To(MatchError(govcd.ErrorEntityNotFound))
		//	Expect(vcdDisk).To(BeNil())
		//
		//})
		//It("PV dynamic provisioning test using delete reclaim policy", func() {
		//	// Todo: create nameSpace4
		//	ns, err := testingClient.CreateNameSpace(context.TODO(), nameSpace4)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(ns).NotTo(BeNil())
		//	sc, err := testingClient.CreateStorageClass(context.TODO(), storageClassDelete, apiv1.PersistentVolumeReclaimDelete, defaultStorageProfile)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(sc).NotTo(BeNil())
		//	//if err != nil {
		//	//	if err != ResourceExisted {
		//	//		Expect(err).NotTo(HaveOccurred())
		//	//	}
		//	//}
		//	By("StorageClass is created successfully")
		//	pvc, err := testingClient.CreatePVC(context.TODO(), nameSpace4, dynamicPVCNameDelete, storageClassDelete, storageSize)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pvc).NotTo(BeNil())
		//	err = testingClient.IsPvcReady(context.TODO(), nameSpace4, dynamicPVCNameDelete)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("PVC is created successfully")
		//	pvc, err = testingClient.GetPVC(context.TODO(), nameSpace4, dynamicPVCNameDelete)
		//	//Todo: Expect
		//	dynamicPVName = pvc.Spec.VolumeName
		//	Expect(dynamicPVCName).NotTo(BeEmpty())
		//	pv, err := testingClient.GetPV(context.TODO(), dynamicPVName)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pv).NotTo(BeNil())
		//	By(fmt.Sprintf("PV [%s] is verified in Kubernetes", dynamicPVName))
		//	vcdDisk, err := VerifyDiskViaVCD(testingClient.VcdClient, dynamicPVName)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(vcdDisk).NotTo(BeNil())
		//	By("PV is verified in VCD")
		//	pvFound, err := GetPVByNameViaRDE(dynamicPVName, testingClient, "named-disk")
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pvFound).To(BeTrue())
		//	By("PV is verified in RDE")
		//	deployment, err := testingClient.CreateDeployment(context.TODO(), &testingsdk.DeployParams{
		//		Name: deployName,
		//		Labels: map[string]string{
		//			"app": "nginx",
		//		},
		//		ContainerParams: testingsdk.ContainerParams{
		//			ContainerName:  "nginx",
		//			ContainerImage: "nginx:1.14.2",
		//			ContainerPort:  80,
		//		},
		//		VolumeParams: testingsdk.VolumeParams{
		//			VolumeName: volumeName,
		//			PvcRef:     dynamicPVCNameDelete,
		//			MountPath:  "/init-container-msg-mount-path",
		//		},
		//	}, nameSpace4)
		//	Expect(deployment).NotTo(BeNil())
		//	Expect(err).NotTo(HaveOccurred())
		//	err = testingClient.IsPvcReady(context.TODO(), nameSpace4, dynamicPVCNameDelete)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("PVC status should be 'bound'")
		//	err = testingClient.IsDeploymentReady(context.TODO(), nameSpace4, dynamicPVCNameDelete)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("Deployment should be ready")
		//	err = testingClient.DeleteDeployment(context.TODO(), nameSpace4, deployName)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("PVC is deleted in Kubernetes")
		//	pv, err = testingClient.GetPV(context.TODO(), dynamicPVName)
		//	Expect(err).To(MatchError(testingsdk.ResourceNotFound))
		//	Expect(pv).To(BeNil())
		//	By("PV is deleted in Kubernetes")
		//	vcdDisk, err = getDiskByNameViaVCD(testingClient.VcdClient, dynamicPVName)
		//	Expect(err).To(MatchError(govcd.ErrorEntityNotFound))
		//	Expect(vcdDisk).To(BeNil())
		//	By("PV is deleted in VCD")
		//	found, err := GetPVByNameViaRDE(dynamicPVName, testingClient, "named-disk")
		//	Expect(err).To(MatchError(testingsdk.ResourceNotFound))
		//	Expect(found).To(BeFalse())
		//	By("PV is deleted in RDE")
		//	err = testingClient.DeleteStorageClass(context.TODO(), storageClassDelete)
		//	Expect(err).NotTo(HaveOccurred())
		//	err = testingClient.DeleteNameSpace(context.TODO(), nameSpace4)
		//	Expect(err).NotTo(HaveOccurred())
		//	//Todo: delete storageClass
		//	//Todo: delete namespace
		//})
		//It("PV dynamic provisioning test using retain reclaim policy", func() {
		//	pvc, err := testingClient.CreatePVC(context.TODO(), nameSpace4, dynamicPVCNameRetain, storageClassRetain, storageSize)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pvc).NotTo(BeNil())
		//	err = testingClient.IsPvcReady(context.TODO(), nameSpace4, dynamicPVCNameRetain)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("PVC is created successfully")
		//	pvc, err = testingClient.GetPVC(context.TODO(), nameSpace4, dynamicPVCNameRetain)
		//	dynamicPVName := pvc.Spec.VolumeName
		//	Expect(dynamicPVCName).NotTo(BeEmpty())
		//	pv, err := testingClient.GetPV(context.TODO(), dynamicPVName)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pv).NotTo(BeNil())
		//	By(fmt.Sprintf("PV [%s] is verified in Kubernetes", dynamicPVName))
		//	vcdDisk, err := getDiskByNameViaVCD(testingClient.VcdClient, dynamicPVName)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(vcdDisk).NotTo(BeNil())
		//	By("PV is verified in VCD")
		//	pvFound, err := GetPVByNameViaRDE(dynamicPVName, testingClient, "named-disk")
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pvFound).To(BeTrue())
		//	By("PV is verified in RDE")
		//	deployment, err := testingClient.CreateDeployment(context.TODO(), &testingsdk.DeployParams{
		//		Name: deployName,
		//		Labels: map[string]string{
		//			"app": "nginx",
		//		},
		//		ContainerParams: testingsdk.ContainerParams{
		//			ContainerName:  "nginx",
		//			ContainerImage: "nginx:1.14.2",
		//			ContainerPort:  80,
		//		},
		//		VolumeParams: testingsdk.VolumeParams{
		//			VolumeName: "nginx-deployment-volume",
		//			PvcRef:     dynamicPVCNameRetain,
		//			MountPath:  "/init-container-msg-mount-path",
		//		},
		//	}, nameSpace4)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(deployment).NotTo(BeNil())
		//
		//	err = testingClient.IsPvcReady(context.TODO(), nameSpace4, dynamicPVCNameRetain)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("PVC status should be 'bound'")
		//	err = testingClient.IsDeploymentReady(context.TODO(), nameSpace4, deployName)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("Deployment should be ready")
		//	err = testingClient.DeleteDeployment(context.TODO(), nameSpace4, deployName)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("Deployment is deleted in Kubernetes")
		//	err = testingClient.DeletePVC(context.TODO(), nameSpace4, dynamicPVCNameRetain)
		//	Expect(err).NotTo(HaveOccurred())
		//	By("PVC is deleted in Kubernetes")
		//	pv, err = testingClient.GetPV(context.TODO(), dynamicPVName)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pv).NotTo(BeNil())
		//	By(fmt.Sprintf("PV [%s] is verified in Kubernetes after PVC is deleted", dynamicPVName))
		//	vcdDisk, err = getDiskByNameViaVCD(testingClient.VcdClient, dynamicPVName)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(vcdDisk).NotTo(BeNil())
		//	By("PV is verified in VCD after PVC is deleted")
		//	pvFound, err = GetPVByNameViaRDE(dynamicPVName, testingClient, "named-disk")
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pvFound).To(BeTrue())
		//	By("PV is verified in RDE after PVC is deleted")
		//
		//	err = testingClient.DeletePV(context.TODO(), dynamicPVName)
		//	Expect(err).NotTo(HaveOccurred())
		//	pv, err = testingClient.GetPV(context.TODO(), dynamicPVName)
		//	Expect(err).To(MatchError(testingsdk.ResourceNotFound))
		//	Expect(pv).To(BeNil())
		//	By("PV is deleted in Kubernetes")
		//	vcdDisk, err = getDiskByNameViaVCD(testingClient.VcdClient, dynamicPVName)
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(vcdDisk).NotTo(BeNil())
		//	By("PV is verified in VCD after PVC is deleted")
		//	pvFound, err = GetPVByNameViaRDE(dynamicPVName, testingClient, "named-disk")
		//	Expect(err).NotTo(HaveOccurred())
		//	Expect(pvFound).To(BeTrue())
		//	By("PV is verified in RDE after PVC is deleted")
		//	err = DeleteDisk(testingClient.VcdClient, dynamicPVName)
		//	Expect(err).NotTo(HaveOccurred())
		//	By(fmt.Sprintf("Disk [%s] is deleted successfully in VCD", dynamicPVName))
		//	err = RemoveDiskViaRDE(testingClient.VcdClient, dynamicPVName, testingClient.ClusterId)
		//	Expect(err).NotTo(HaveOccurred())
		//	By(fmt.Sprintf("Disk [%s] is deleted successfully in RDE", dynamicPVName))
		//	vcdDisk, err = getDiskByNameViaVCD(testingClient.VcdClient, dynamicPVName)
		//	Expect(err).To(MatchError(govcd.ErrorEntityNotFound))
		//	Expect(vcdDisk).To(BeNil())
		//	By(fmt.Sprintf("Disk [%s] is not shown after deletion in RDE", dynamicPVName))
		//	found, err := GetPVByNameViaRDE(dynamicPVName, testingClient, "named-disk")
		//	Expect(err).To(MatchError(testingsdk.ResourceNotFound))
		//	Expect(found).To(BeFalse())
		//	By(fmt.Sprintf("Disk [%s] is not shown after deletion in RDE", dynamicPVName))
		//})

		AfterEach(func() {
			err := testingClient.DeleteStorageClass(context.TODO(), storageClassRetain)
			Expect(err).NotTo(HaveOccurred())
			err = testingClient.DeleteStorageClass(context.TODO(), storageClassDelete)
			Expect(err).NotTo(HaveOccurred())
			err = testingClient.DeleteNameSpace(context.TODO(), nameSpace4)
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
