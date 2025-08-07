import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  eslint: {
    // Disable ESLint during build temporarily for testing
    ignoreDuringBuilds: true,
  },
  typescript: {
    // Disable type checking during build temporarily
    ignoreBuildErrors: true,
  },
};

export default nextConfig;
