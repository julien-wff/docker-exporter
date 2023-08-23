import 'dotenv/config';
import Docker from 'dockerode';
import fastify from 'fastify';
import promClient from 'prom-client';
import { dockerMetricsToRegistry, registerDockerMetrics } from './lib/dockerMetricsToRegistry';
import { getDockerMetrics } from './lib/getDockerMetrics';

// Setup web server
const server = fastify({
    logger: process.env.LOGGER_ENABLED !== 'false' && process.env.LOGGER_ENABLED !== '0',
});

// Setup Docker client
const docker = new Docker();

// Setup prometheus client
const register = new promClient.Registry();
registerDockerMetrics(register);

// Add metric to display the duration of the last scrape
const lastScrapeDuration = new promClient.Gauge({
    name: 'docker_last_scrape_duration_seconds',
    help: 'Duration of the last scrape of metrics from Docker in seconds',
    registers: [ register ],
});
lastScrapeDuration.set(0);

// Setup prometheus metrics endpoint
server.get('/metrics', async (_, res) => {
    const start = Date.now();
    const metrics = await getDockerMetrics(docker);
    dockerMetricsToRegistry(metrics);
    lastScrapeDuration.set((Date.now() - start) / 1000);
    res.header('Content-Type', register.contentType);
    res.send(await register.metrics());
});

// Setup docker healthcheck endpoint
server.get('/health', async (_, res) => {
    try {
        await docker.ping();
        res.status(200).send('OK');
    } catch (err) {
        res.status(500).send((err as Error).message);
    }
});

// Get port and host from environment variables
const port = isNaN(parseInt(process.env.SERVER_PORT as string)) ? 9100 : parseInt(process.env.SERVER_PORT as string);
const host = process.env.SERVER_HOST as string | undefined || '0.0.0.0';

// Start server
server.listen({ port, host }, err => {
    if (err) {
        console.error(err);
        process.exit(1);
    }
});
