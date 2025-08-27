#!/bin/bash
echo "Rolling back to BEFORE_CLEANUP..."
git reset --hard BEFORE_CLEANUP_$(date +%Y%m%d)
echo "Rollback complete"