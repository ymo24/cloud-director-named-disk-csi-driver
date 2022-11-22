package vcdcsiclient

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"testing"
)

//const kubeconfig = "apiVersion: v1\\nclusters:\\n- cluster:\\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM2akNDQWRLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJeU1URXhOekEwTlRReU5Wb1hEVE15TVRFeE5EQTBOVGt5TlZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBSzRaCldqQkxQeHRWYzB0V3lrWnd5am85WUJRT2RnR3p2SFd2UFVuWmorMzFDK0Jac2xFSXFyS3dJQVpZL3FVdEZiUzYKcFdtNHhnVVYxNExzZExtNTFGUTFNOU9PYjFTZzlYWG1nS2pDWTEvem5pZlhCSGRLWmUvK2JjdjNUTWJqZGxWUApKZkRSQW1WMnpEMWJQU2FNVE1lK2hpMkNoT3RkS0JBbmRNNWg1cG1UZnlVcTNrcTRjNU13dnpITDJzY3R4NGx1CnJxK0NHSFlvV2VlKzlKVlZ1c2J6aWZ3UzZIR1RvTUZZTlVDcHF4ZDFzc290Z08wN1BJOEhPaTdBNFFFdU1zL1QKdVFzUk84dVl4bEo2V3hVbDVicUkyNklPU204UTdVRXZRMVFYLzJ1Kzh1c3JPUXk4ZmJvQ3UvcXNNN0lGUmN3Nwp4ZlZRemZTWjR6VHp6dDhQT2VVQ0F3RUFBYU5GTUVNd0RnWURWUjBQQVFIL0JBUURBZ0trTUJJR0ExVWRFd0VCCi93UUlNQVlCQWY4Q0FRQXdIUVlEVlIwT0JCWUVGTlFOMGxaL0h1eXJ3VTJ1c2EyQXpXajdXK29MTUEwR0NTcUcKU0liM0RRRUJDd1VBQTRJQkFRQWh5VmVBN3BodVJHQUxHMmQ4WGx2bHFlV04xaHU2ZzJ4MW4yQ0RoUFpMZnFlTApTNXlVWC82RjkrRDFWbmpDeUZWMnNUYVVBeDkrZEVsV0ExYmEyOXAwYjJEekhPLzB1S1hOS01nd0plL2pMTFY1CnZxSmkyUHVTNmN3WkswN0E1dERqdWVONURxUEVOT1NOV0c3WDhKS0F1WVVrR3FaUjZma1JwYWpiM0ZOTmFXQzEKQTBPZHFzOTJDSHZGdG5yMGZaUzJQaDlWeFdON3N4SnduZ3k4REJLTjJuemNyMzlKdHYzUkNya203MGJoMlpmYwowSysrSkpXR2FMZXB2OUpRTVRuY1QrUnJ5L21nUFV2ZEtIKzlUWis5alMya1FGdXRnSFJrMFE0THB2ZXZVNU9LCkhOcWhYdXVUTWxRODFxRmwyanZ0VStNbUpGNTErUi8zUlJHWGpZb0gKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=\\n    server: https://10.89.104.62:6443\\n  name: official-cluster\\ncontexts:\\n- context:\\n    cluster: official-cluster\\n    user: official-cluster-admin\\n  name: official-cluster-admin@official-cluster\\ncurrent-context: official-cluster-admin@official-cluster\\nkind: Config\\npreferences: {}\\nusers:\\n- name: official-cluster-admin\\n  user:\\n    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURFekNDQWZ1Z0F3SUJBZ0lJVjFqUmJ2SkJjVXd3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TWpFeE1UY3dORFUwTWpWYUZ3MHlNekV4TVRjd05EVTVNamRhTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXZGSk1GL3IraENhME5sL2UKWEplc2N3ajQ1amtjSHZ5RmV3ZUI3eTY5L3FLTTFIQXVyb0U5VEM2MHFzNGE3cnN5UDZzSDRuY0ZRdmdPTWxVTQpmODA1ZXJTaTBGR2N0M3VHWnZWd3NxZDIzSVRaLzd4dWtxR3hBQXh5Mkh6b2pCRkxXb01FTnRONDV1OS9LM3RlCnlKdm5Vblh2VlBDZnA1ZHdUWUFzNU13MmFCZ3N1M0R3WVZOdUFwV3VHbFF1dDFoK2lFYkNRWmhMTFArUDJhTUgKdHduZ2N6TXFzeVVSRzNkRzNTUUJPWlhDQmp2eGZNWW1vWTVvUzIrb09QUkhKbDFSM053YWdZRk82V2JHekFubgpmekZoWU9VUkVtd2F0NmlHYXBZUnNVemM1M1RyUjYvSm9Zb21LWW11azlLODBVbGc5MWVqVUxlWG9zbjlMZlJmCmhRMGVpd0lEQVFBQm8wZ3dSakFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0h3WURWUjBqQkJnd0ZvQVUxQTNTVm44ZTdLdkJUYTZ4cllETmFQdGI2Z3N3RFFZSktvWklodmNOQVFFTApCUUFEZ2dFQkFBWmJvOFF2RW5VRUIzYWF5bkFQY00xcy82VHVERDNqa2paRFUrenNQRkxlb0gzYjg2VDNEUG5xCldNUmtYMGtXem14eGIrQVEybkhDRkh5MWx6Z1dOQnQwa2VPcFNaeVRRR3hlQmJJMm1ONFZEZjdRRklUcXRmYngKTnNocVRqLzhObUxoWW9nWEN2akRmQUphbVp4RlMybXZPQjlKNEF6RW5pcktuakxNK0RZMHFmOWF2VDhjYmQ0UQo1UVpScStVeGl4OW1RU0NrNWNGc3V3aFI1d2ZrN3NieFVEMUZPUWhVbUJLQ1Y0MDBrUUZzRXcxbHdoZHJ5cFM3CldMRVNuMXB5WFU1bm9YM2VrbTNuNUVkakNwVncwTlB4ZkFMbFZkRHQrOFhUcG1KbE5HaktjYzBod2dySEkzWkUKSDVpdkkyc0RmRmhFcjFlOXVpYUV6QnpzSkpZZUFmTT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=\\n    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb2dJQkFBS0NBUUVBdkZKTUYvcitoQ2EwTmwvZVhKZXNjd2o0NWprY0h2eUZld2VCN3k2OS9xS00xSEF1CnJvRTlUQzYwcXM0YTdyc3lQNnNING5jRlF2Z09NbFVNZjgwNWVyU2kwRkdjdDN1R1p2VndzcWQyM0lUWi83eHUKa3FHeEFBeHkySHpvakJGTFdvTUVOdE40NXU5L0szdGV5SnZuVW5YdlZQQ2ZwNWR3VFlBczVNdzJhQmdzdTNEdwpZVk51QXBXdUdsUXV0MWgraUViQ1FaaExMUCtQMmFNSHR3bmdjek1xc3lVUkczZEczU1FCT1pYQ0JqdnhmTVltCm9ZNW9TMitvT1BSSEpsMVIzTndhZ1lGTzZXYkd6QW5uZnpGaFlPVVJFbXdhdDZpR2FwWVJzVXpjNTNUclI2L0oKb1lvbUtZbXVrOUs4MFVsZzkxZWpVTGVYb3NuOUxmUmZoUTBlaXdJREFRQUJBb0lCQURMRGRHT2FjdlVvZ2JlTwpqQ0FsZW9UZnpFZ3k2Tk9wZWttNXNsckpITW9CQXpsWXJCeTZGYzN0WVNJUndNek5oVFFJWkcyMWE2T3J2aGZ2Ck9SbFNOc1pPM3Z5TW8xUUtaaVptenVRNXBCNjVhUkk0dHcycnJFeTVEbFF4QjNhS1N1ZXhIWGV4OVlzNnorcm0Kckp3aSttNE9BRi9ESlNaRitpM0orVkFMaERiMnBwSVRDYUtzNDE3LzRpTytyMFlSeWNnYVRSSmFKTFF3MTlQWApPSi9FTEs2dFRBVXlPRVFLVDlZRk1MOWdMN2xPankvSm8xTVZEeTh1SXNST2o4VHpyVmZBZDNGUWFWTmVCeEV4CjZJaTlsV3BtVWJxRjJRUTRYK0puSHZ6dXI3Ulp5MVNwRDNoQ0txdkovWlVuTDlVallUWmZweWNKcTMwTzVDc04KOTJFTitURUNnWUVBeTUwaHR1QnJFZnRRcnl6Q1VwTklsNHp6dGxKdVV2K0JwNTFRRVJaTnA0SEVPb01mYlNJNApSYVRVWnpCTGRuVml2WnJ4L3lEZTFDSTcxQ2owV2VmWXFCK0VnOENHMVl3SktQeUtuTjNScDZoOEorbjh2S21mCjNDL29iSk9rTUZDeURCRkFrZFZnbnJ2YzhuZkpxaS80RzV1OEdBczlhcFRQMVcvUENtSkN2elVDZ1lFQTdNWHkKeG9rRDNPbjZzcWIxQmJnc3hVNVVDMElYN1Z6TUlVTkR2YkZGdy96TGNId1l6am5XME1IQ1BvZjhxak03VkU3bgo1M0VIcHVGSFJyV3phWjB1bzdnM0RiTTJ4dDNMbVFMQXJlb0FiZjF6QmhIc0J0RTg2eHZ4VEY4YS8vU3J2NnhpCjdYZk0yLy9MSlhGK0xQaGtPR0NiNWhBTkZFc1E4a3pobjdFL1hyOENnWUFVb2ZLVnBuNFREMlZvcXQ5eUlLeWQKZHRJSGFxajFUaURrVEVPZHg2WE0wSkNDNFdDZzNYUFlVdENYT0VTZFYxM1BHdEZrNmY3S2ZrR1R5U2FocWFYQgp1NWZoQmZSajFWSGtUbHI1ZEZ6WFlYSmJWUkdnU0l3RGN2TlpkVWlSQU14Wi9yR05WWkw1NHMyTDRHbVdEbEJVCjg5NEdqYlVHaE5mZXAvclI4WTBUOVFLQmdGeEw4WTVZM011aDNkc0VZQ0Vob1Rvc2hYQjZEQyszNjg3UGxMbCsKUUE3ZEhVUzA2MHFBbTI5M0NFd3Q0RjFNYVVVOUdRTk1PVXBoS05LMGc4S1l4aFNGKzlmNFUxTUVKSmg1elRnKwpMbnF4d01QTitxN1JvNmlXbE9KTGRJL3dCWDlMS0trZStSbU5SZGhMdkg2MU9RUU5ETmlLeXo1czRLZEROdlIwCk9KcXJBb0dBY2xmRElEWFdvVlo2bmZuL1VNTVl1ZnRLYjNxQnR0WEp4TjJObEdmSkVTQUJqOHYweWNkRm9RSjAKVFhRdU45a2ZnREtRUnczK1R5WDkxOEprNzY1dWFLZzN5NWt0K3ROblVockRLcFM3bkZvU2pMUThVOCtrN05YVQpmSE1sa1BRMEVCYk1Eak8zN29KRjlwRDhRdURRY3liNGNYNnJKYXdpR3U5WEZKQUtMTFU9Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==\\n"
const (
	kubeconfig = "/Users/ymo/Downloads/kubeconfig-official-cluster.txt"
)

type Book struct {
	Title  string
	Author string
	Pages  int
}

func (b *Book) CategoryByLength() string {

	if b.Pages >= 300 {
		return "NOVEL"
	}

	return "SHORT STORY"
}

func TestCSIAutomation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CSI Testing Suite")
}

var _ = Describe("CSI Driver Automation", func() {
	var (
		clientset *kubernetes.Clientset
		csiDriver *DiskManager
	)

	BeforeEach(func() {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		csiDiskManager := new(DiskManager)
		cloudConfig, err := getTestConfig()
		vcdClient, err := getTestVCDClient(cloudConfig, map[string]interface{}{
			"getVdcClient": true,
			"user":         "administrator",
			"refreshToken": "3X7P91KUzcbpv1JIeJlJ7AwHwvTeUq7E",
			"userOrg":      "org1",
		})
		csiDiskManager.VCDClient = vcdClient
		csiDriver = csiDiskManager
		csiDriver.ClusterID = cloudConfig.ClusterID
	})

	//Describe("CSI Driver Automation", func() {
	//
	//	Context("With fewer than 300 pages", func() {
	//		It("should be a short story", func() {
	//			//Expect(shortBook.CategoryByLength()).To(Equal("SHORT STORY"))
	//
	//			pods, err := clientset.CoreV1().Pods("kube-system").List(metav1.ListOptions{})
	//			//if err != nil {
	//			if err != nil {
	//				panic(err.Error())
	//			}
	//			fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	//		})
	//	})
	//	Context("Persistent Volume Creation Test", func() {
	//		It("Persistent Volume should be created successfully", func() {
	//			_, err := csiDriver.CreateDisk("test-disk-02", 100, VCDBusTypeSCSI, VCDBusSubTypeVirtualSCSI,
	//				"", "*", true)
	//			if err != nil {
	//				panic(err.Error())
	//			}
	//		})
	//	})
	//
	//	Context("Deployment Installation Test", func() {
	//		It("Persistent Volume should be created successfully", func() {
	//			deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	//
	//			deployment := &appsv1.Deployment{
	//				ObjectMeta: metav1.ObjectMeta{
	//					Name: "wordpress-mysql",
	//				},
	//				Spec: appsv1.DeploymentSpec{
	//					Selector: &metav1.LabelSelector{
	//						MatchLabels: map[string]string{
	//							"app": "demo",
	//						},
	//					},
	//					Template: apiv1.PodTemplateSpec{
	//						ObjectMeta: metav1.ObjectMeta{
	//							Labels: map[string]string{
	//								"app": "demo",
	//							},
	//						},
	//						Spec: apiv1.PodSpec{
	//							Containers: []apiv1.Container{
	//								{
	//									Name:  "web",
	//									Image: "nginx:1.12",
	//									Ports: []apiv1.ContainerPort{
	//										{
	//											Name:          "http",
	//											Protocol:      apiv1.ProtocolTCP,
	//											ContainerPort: 80,
	//										},
	//									},
	//									VolumeMounts: []apiv1.VolumeMount{
	//										{
	//											Name:      "",
	//											MountPath: "",
	//										},
	//									},
	//								},
	//							},
	//							Volumes: []apiv1.Volume{},
	//						},
	//					},
	//				},
	//			}
	//
	//		})
	//	})
	//
	//})

	It("should create persistent Volume", func() {
		By("Input PV Param")
		pv := &apiv1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-pv-01",
			},
			Spec: apiv1.PersistentVolumeSpec{
				StorageClassName: "default-class-1",
				AccessModes: []apiv1.PersistentVolumeAccessMode{
					"ReadWriteOnce",
				},
				PersistentVolumeSource: apiv1.PersistentVolumeSource{
					//NFS: &apiv1.NFSVolumeSource{
					//	Path:   "/tmp",
					//	Server: "172.17.0.2",
					//},
					CSI: &apiv1.CSIPersistentVolumeSource{
						Driver:       "named-disk.csi.cloud-director.vmware.com",
						FSType:       "ext4",
						VolumeHandle: "existingVolumeName",
					},
				},
				Capacity: apiv1.ResourceList{
					"storage": resource.MustParse("2Gi"),
				},
				ClaimRef: &apiv1.ObjectReference{
					Namespace: apiv1.NamespaceDefault,
					Name:      "mysql-pv-claim",
				},
			},
		}
		_, err := clientset.CoreV1().PersistentVolumes().Create(pv)
		Expect(err).NotTo(HaveOccurred())
		//existingPVs, err := clientset.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
		_, err = clientset.CoreV1().PersistentVolumes().Get("test-pv", metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())
		// Todo: check PV in VCD
	})

	//Describe("Deployment Installation", func() {
	//	It("should create a successful deployment", func() {
	//		// Todo: Add PVC
	//		deploymentClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	//
	//		deployment := &appsv1.Deployment{
	//			ObjectMeta: metav1.ObjectMeta{
	//				Name: "wordpress-mysql",
	//			},
	//			Spec: appsv1.DeploymentSpec{
	//				Selector: &metav1.LabelSelector{
	//					MatchLabels: map[string]string{
	//						"app":  "wordpress",
	//						"tier": "mysql",
	//					},
	//				},
	//				Strategy: appsv1.DeploymentStrategy{
	//					Type: appsv1.RecreateDeploymentStrategyType,
	//				},
	//				Template: apiv1.PodTemplateSpec{
	//					ObjectMeta: metav1.ObjectMeta{
	//						Labels: map[string]string{
	//							"app":  "wordpress",
	//							"tier": "mysql",
	//						},
	//					},
	//					Spec: apiv1.PodSpec{
	//						Containers: []apiv1.Container{
	//							{
	//								Name:  "mysql",
	//								Image: "mysql:5.6",
	//								Env: []apiv1.EnvVar{
	//									{
	//										Name: "MYSQL_ROOT_PASSWORD",
	//										ValueFrom: &apiv1.EnvVarSource{
	//											SecretKeyRef: &apiv1.SecretKeySelector{
	//												Key: "password",
	//												LocalObjectReference: apiv1.LocalObjectReference{
	//													Name: "mysql-pass",
	//												},
	//											},
	//										},
	//									},
	//								},
	//								Ports: []apiv1.ContainerPort{
	//									{
	//										Name:          "http",
	//										Protocol:      apiv1.ProtocolTCP,
	//										ContainerPort: 3306,
	//									},
	//								},
	//								VolumeMounts: []apiv1.VolumeMount{
	//									{
	//										Name:      "mysql-persistent-storage",
	//										MountPath: "/var/lib/mysql",
	//									},
	//								},
	//							},
	//						},
	//						Volumes: []apiv1.Volume{
	//							{
	//								Name: "mysql-persistent-storage",
	//								VolumeSource: apiv1.VolumeSource{
	//									PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
	//										ClaimName: "mysql-pv-claim",
	//									},
	//								},
	//							},
	//						},
	//					},
	//				},
	//			},
	//		}
	//
	//		result, err := deploymentClient.Create(deployment)
	//		Expect(err).NotTo(HaveOccurred())
	//		Expect(result).NotTo(BeNil())
	//
	//	})
	//})
	//
	//Describe("PVC Verification in Kubernetes", func() {
	//	It("PVC should be visible", func() {
	//		pvc, err := clientset.CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault).Get("mysql-pv-claim", metav1.GetOptions{})
	//		Expect(err).NotTo(HaveOccurred())
	//		Expect(pvc).NotTo(BeNil())
	//		By("PVC status should be 'bound'")
	//	})
	//})
})
