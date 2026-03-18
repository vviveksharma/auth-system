# Core Static Operations
POST   /api/v1/statics/upload              # Upload new static
GET    /api/v1/statics/{id}                 # Get metadata
GET    /api/v1/statics/{id}/download        # Download file
DELETE /api/v1/statics/{id}                  # Delete static
PATCH  /api/v1/statics/{id}                  # Update metadata

# Project-scoped operations
GET    /api/v1/projects/{projectId}/statics  # List all statics in project
POST   /api/v1/projects/{projectId}/statics  # Upload to specific project

# Batch operations
POST   /api/v1/statics/batch/delete          # Bulk delete
POST   /api/v1/statics/batch/metadata        # Bulk metadata update

# Signed URLs for private assets
POST   /api/v1/statics/{id}/signed-url       # Generate temporary access URL


##### Monthly usage 

GET    /api/v1/stats/requests               # List requests (paginated)
GET    /api/v1/stats/requests/summary        # Aggregated summary
GET    /api/v1/stats/requests/top-endpoints  # Most called endpoints
GET    /api/v1/stats/requests/errors         # Failed requests
GET    /api/v1/projects/{projectId}/requests # Project-specific requests