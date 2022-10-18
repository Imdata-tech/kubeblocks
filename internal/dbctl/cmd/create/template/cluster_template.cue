// required, command line input options for parameters and flags
options: {
	name:          string
	namespace:     string
	clusterDefRef: string
	appVersionRef: string
	components: [...]
	terminationPolicy: string
	podAntiAffinity:   string
	topologyKeys: [...]
	nodeLabels: {}
}

// required, k8s api resource content
content: {
	apiVersion: "dbaas.infracreate.com/v1alpha1"
	kind:       "Cluster"
	metadata: {
		name:      options.name
		namespace: options.namespace
	}
	spec: {
		clusterDefinitionRef: options.clusterDefRef
		appVersionRef:        options.appVersionRef
		affinity: {
			podAntiAffinity: options.podAntiAffinity
			topologyKeys:    options.topologyKeys
			nodeLabels:      options.nodeLabels
		}
		components:        options.components
		terminationPolicy: options.terminationPolicy
	}
}
