import { Gauge, type Registry } from 'prom-client';
import { type DockerMetrics } from './getDockerMetrics';

const dockerImageSize = new Gauge({
    name: 'docker_image_size_bytes',
    help: 'The size of the image in bytes',
    labelNames: [ 'id', 'tag', 'containers' ] as const,
});

const dockerImageVirtualSize = new Gauge({
    name: 'docker_image_virtual_size_bytes',
    help: 'The virtual size of the image in bytes',
    labelNames: [ 'id', 'tag', 'containers' ] as const,
});

const dockerVolumeSize = new Gauge({
    name: 'docker_volume_size_bytes',
    help: 'The size of the volume in bytes',
    labelNames: [ 'name', 'mountpoint', 'compose_project', 'containers' ] as const,
});

const dockerNetworkContainerCount = new Gauge({
    name: 'docker_network_container_count',
    help: 'The number of containers connected to the network',
    labelNames: [ 'id', 'name' ] as const,
});

export function registerDockerMetrics(register: Registry) {
    register.registerMetric(dockerImageSize);
    register.registerMetric(dockerImageVirtualSize);
    register.registerMetric(dockerVolumeSize);
    register.registerMetric(dockerNetworkContainerCount);
}

export function dockerMetricsToRegistry(metrics: DockerMetrics) {
    for (const image of metrics.images) {
        dockerImageSize.set({
            id: image.id,
            tag: image.repoTags[0] ?? '<none>',
            containers: image.containerCount,
        }, image.size);
        dockerImageVirtualSize.set({
            id: image.id,
            tag: image.repoTags[0] ?? '<none>',
            containers: image.containerCount,
        }, image.virtualSize);
    }

    for (const volume of metrics.volumes) {
        dockerVolumeSize.set({
            name: volume.name,
            mountpoint: volume.mountPoint,
            compose_project: volume.composeProject || '',
            containers: volume.containerCount,
        }, volume.size);
    }

    for (const network of metrics.networks) {
        dockerNetworkContainerCount.set({ id: network.id, name: network.name }, network.containerCount);
    }
}
