import type Docker from 'dockerode';
import ffs from 'fast-folder-size';
import path from 'node:path';
import util from 'node:util';

const fastFolderSize = util.promisify(ffs);

export async function getDockerMetrics(docker: Docker): Promise<DockerMetrics> {

    const images = await docker.listImages({ all: true });
    const volumes = await docker.listVolumes();
    const networks = await docker.listNetworks();

    const result = {
        images: images.map(image => ({
            id: image.Id,
            repoTags: image.RepoTags || [],
            size: image.Size,
            virtualSize: image.VirtualSize,
            containerCount: 0,
        }) satisfies DockerImageMetrics),
        volumes: volumes.Volumes.map(volume => ({
            name: volume.Name,
            mountPoint: volume.Mountpoint,
            size: volume.UsageData?.Size || 0,
            composeProject: volume.Labels?.['com.docker.compose.project'] || null,
            containerCount: 0,
        }) satisfies DockerVolumeMetrics),
        networks: networks.map(network => ({
            id: network.Id,
            name: network.Name,
            containerCount: 0,
        }) satisfies DockerNetworkMetrics),
    };

    for (const container of await docker.listContainers({ all: true })) {
        for (const image of result.images) {
            if (container.ImageID === image.id) {
                image.containerCount++;
            }
        }

        for (const volume of result.volumes) {
            if (container.Mounts?.some(mount => mount.Name === volume.name)) {
                volume.containerCount++;
            }
        }

        for (const network of result.networks) {
            if (container.NetworkSettings?.Networks?.[network.name]) {
                network.containerCount++;
            }
        }
    }

    for (const volume of result.volumes) {
        // Get the size of the volume from the host
        try {
            const dir = path.join('/rootfs', volume.mountPoint);
            const stats = await fastFolderSize(dir);
            volume.size = stats || 0;
        } catch (err) {
            // Ignore errors
        }
    }

    return result;
}

export interface DockerMetrics {
    images: DockerImageMetrics[];
    volumes: DockerVolumeMetrics[];
    networks: DockerNetworkMetrics[];
}

export interface DockerImageMetrics {
    id: string;
    repoTags: string[];
    size: number;
    virtualSize: number;
    containerCount: number;
}

export interface DockerVolumeMetrics {
    name: string;
    mountPoint: string;
    size: number;
    composeProject: string | null;
    containerCount: number;
}

export interface DockerNetworkMetrics {
    id: string;
    name: string;
    containerCount: number;
}
