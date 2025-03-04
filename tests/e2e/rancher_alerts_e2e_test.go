package e2e_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rancher/observability-e2e/tests/helper/utils"
	rancher "github.com/rancher/shepherd/clients/rancher"
	"github.com/rancher/shepherd/extensions/kubectl"
	e2e "k8s.io/kubernetes/test/e2e/framework"
)

// Define the struct for the Alert
type RancherAlert struct {
	Annotations  map[string]string  `json:"annotations"`
	EndsAt       string             `json:"endsAt"`
	Fingerprint  string             `json:"fingerprint"`
	Receivers    []RancherReceiver  `json:"receivers"`
	StartsAt     string             `json:"startsAt"`
	Status       RancherAlertStatus `json:"status"`
	UpdatedAt    string             `json:"updatedAt,omitempty"`
	GeneratorURL string             `json:"generatorURL,omitempty"`
	Labels       map[string]string  `json:"labels"`
}

// Define the struct for Receiver
type RancherReceiver struct {
	Name string `json:"name"`
}

// Define the struct for AlertStatus
type RancherAlertStatus struct {
	InhibitedBy []string `json:"inhibitedBy"`
	SilencedBy  []string `json:"silencedBy"`
	State       string   `json:"state"`
}

const (
	alertmanagerConfigFilePath = "../helper/yamls/alertManagerConfig.yaml"
)

var _ = Describe("Observability Alerting E2E Test Suite", func() {
	var clientWithSession *rancher.Client // RancherConfig *Config

	JustBeforeEach(func() {
		By("Creating a client session")
		clientWithSession, err = client.WithSession(sess)
		Expect(err).NotTo(HaveOccurred())
	})

	It("[QASE-6831] Test : Verify status of rancher-alert Deployments using kubectl", Label("LEVEL1", "alerts", "E2E"), func() {
		testCaseID = 6831
		By("1) Fetch all the deployments belonging to rancher-alerts")
		fetchDeployments := []string{"kubectl", "get", "deployments", "-n", "cattle-monitoring-system", "--no-headers"}
		rancherAlertsDeployments, err := kubectl.Command(clientWithSession, nil, "local", fetchDeployments, "")
		Expect(err).NotTo(HaveOccurred(), "Failed to get deployments")

		By("2) Read all the deployments and verify the status of rancher-alerts deployments")
		foundRancherAlerting := false
		deployments := strings.Split(rancherAlertsDeployments, "\n")
		for _, deployment := range deployments {
			if deployment == "" {
				continue
			}

			fields := strings.Fields(deployment)
			Expect(len(fields)).To(BeNumerically(">=", 4), "Unexpected output format for deployment: %s", deployment)

			deploymentName := fields[0]
			readyReplicas := fields[1]
			availableReplicas := fields[3]

			readyCount := strings.Split(readyReplicas, "/")[0]
			desiredCount := strings.Split(readyReplicas, "/")[1]

			if strings.HasPrefix(deploymentName, "rancher-alerting") {
				foundRancherAlerting = true
				Expect(availableReplicas).To(Equal(desiredCount), "Deployment %s is not fully available. Desired: %s, Available: %s", deploymentName, desiredCount, availableReplicas)
				Expect(readyCount).To(Equal(desiredCount), "Deployment %s pods are not fully ready. Desired: %s, Ready: %s", deploymentName, desiredCount, readyCount)
			}
		}
		Expect(foundRancherAlerting).To(BeTrue(), "No deployments found starting with 'rancher-alerting'")
	})

	It("[QASE-6832] Test : Verify status of rancher-alerts pods using kubectl", Label("LEVEL1", "alerts", "E2E"), func() {
		testCaseID = 6832
		By("1) Fetch all the pods belongs to rancher-alerts")
		fetchPods := []string{"kubectl", "get", "pods", "-n", "cattle-monitoring-system", "--no-headers"}
		rancherAlertsPods, err := kubectl.Command(clientWithSession, nil, "local", fetchPods, "")
		Expect(err).NotTo(HaveOccurred(), "Failed to get pods")

		By("2) Read all the pods and verify the status of rancher-alerts-Pods")
		rancherAlertingFoundPod := false
		alertmanagerFoundPod := false
		pods := strings.Split(rancherAlertsPods, "\n")
		for _, pod := range pods {
			if pod == "" {
				continue
			}
			fields := strings.Fields(pod) // Split the line into pod name and its current status
			Expect(len(fields)).To(BeNumerically(">=", 3), "Unexpected output format for pod: %s", pod)

			podName := fields[0]
			podStatus := fields[2]

			if strings.HasPrefix(podName, "rancher-alerting") && podStatus == "Running" {
				rancherAlertingFoundPod = true
			}
			if strings.HasPrefix(podName, "alertmanager") && podStatus == "Running" {
				alertmanagerFoundPod = true
			}
		}

		Expect(rancherAlertingFoundPod).To(BeTrue(), "Pod with name 'rancher-alerting' is not running or not present")
		Expect(alertmanagerFoundPod).To(BeTrue(), "Pod with name 'alertmanager' is not running or not present")
	})

	It("[QASE-6833] Test : Verify Creating alert manager config using kubectl", Label("LEVEL1", "alerts", "E2E", "AMC"), func() {
		testCaseID = 6833
		By("1) Apply yaml to create alert manager config")
		alertManagerConfigError := utils.DeployAlertManagerConfig(clientWithSession, alertmanagerConfigFilePath)
		if alertManagerConfigError != nil {
			e2e.Logf("Failed to deploy AMC rule: %v", alertManagerConfigError)
		}

		By("2) Fetch all the AMC")
		fetchAlertManagerConfig := []string{"kubectl", "get", "AlertmanagerConfig", "amc", "-n", "cattle-monitoring-system"}
		verifyAlertManagerConfig, err := kubectl.Command(clientWithSession, nil, "local", fetchAlertManagerConfig, "")
		Expect(err).NotTo(HaveOccurred(), "Failed to fetch alert manager config 'amc'")

		e2e.Logf("Successfully fetched AMC: %v", verifyAlertManagerConfig)
	})

})
