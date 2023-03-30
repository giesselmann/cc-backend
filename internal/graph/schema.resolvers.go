package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ClusterCockpit/cc-backend/internal/auth"
	"github.com/ClusterCockpit/cc-backend/internal/graph/generated"
	"github.com/ClusterCockpit/cc-backend/internal/graph/model"
	"github.com/ClusterCockpit/cc-backend/internal/metricdata"
	"github.com/ClusterCockpit/cc-backend/internal/repository"
	"github.com/ClusterCockpit/cc-backend/pkg/archive"
	"github.com/ClusterCockpit/cc-backend/pkg/schema"
)

// Partitions is the resolver for the partitions field.
func (r *clusterResolver) Partitions(ctx context.Context, obj *schema.Cluster) ([]string, error) {
	return r.Repo.Partitions(obj.Name)
}

// Tags is the resolver for the tags field.
func (r *jobResolver) Tags(ctx context.Context, obj *schema.Job) ([]*schema.Tag, error) {
	return r.Repo.GetTags(&obj.ID)
}

// MetaData is the resolver for the metaData field.
func (r *jobResolver) MetaData(ctx context.Context, obj *schema.Job) (interface{}, error) {
	return r.Repo.FetchMetadata(obj)
}

// UserData is the resolver for the userData field.
func (r *jobResolver) UserData(ctx context.Context, obj *schema.Job) (*model.User, error) {
	return auth.FetchUser(ctx, r.DB, obj.User)
}

// CreateTag is the resolver for the createTag field.
func (r *mutationResolver) CreateTag(ctx context.Context, typeArg string, name string) (*schema.Tag, error) {
	id, err := r.Repo.CreateTag(typeArg, name)
	if err != nil {
		return nil, err
	}

	return &schema.Tag{ID: id, Type: typeArg, Name: name}, nil
}

// DeleteTag is the resolver for the deleteTag field.
func (r *mutationResolver) DeleteTag(ctx context.Context, id string) (string, error) {
	panic(fmt.Errorf("not implemented: DeleteTag - deleteTag"))
}

// AddTagsToJob is the resolver for the addTagsToJob field.
func (r *mutationResolver) AddTagsToJob(ctx context.Context, job string, tagIds []string) ([]*schema.Tag, error) {
	jid, err := strconv.ParseInt(job, 10, 64)
	if err != nil {
		return nil, err
	}

	tags := []*schema.Tag{}
	for _, tagId := range tagIds {
		tid, err := strconv.ParseInt(tagId, 10, 64)
		if err != nil {
			return nil, err
		}

		if tags, err = r.Repo.AddTag(jid, tid); err != nil {
			return nil, err
		}
	}

	return tags, nil
}

// RemoveTagsFromJob is the resolver for the removeTagsFromJob field.
func (r *mutationResolver) RemoveTagsFromJob(ctx context.Context, job string, tagIds []string) ([]*schema.Tag, error) {
	jid, err := strconv.ParseInt(job, 10, 64)
	if err != nil {
		return nil, err
	}

	tags := []*schema.Tag{}
	for _, tagId := range tagIds {
		tid, err := strconv.ParseInt(tagId, 10, 64)
		if err != nil {
			return nil, err
		}

		if tags, err = r.Repo.RemoveTag(jid, tid); err != nil {
			return nil, err
		}
	}

	return tags, nil
}

// UpdateConfiguration is the resolver for the updateConfiguration field.
func (r *mutationResolver) UpdateConfiguration(ctx context.Context, name string, value string) (*string, error) {
	if err := repository.GetUserCfgRepo().UpdateConfig(name, value, auth.GetUser(ctx)); err != nil {
		return nil, err
	}

	return nil, nil
}

// Clusters is the resolver for the clusters field.
func (r *queryResolver) Clusters(ctx context.Context) ([]*schema.Cluster, error) {
	return archive.Clusters, nil
}

// Tags is the resolver for the tags field.
func (r *queryResolver) Tags(ctx context.Context) ([]*schema.Tag, error) {
	return r.Repo.GetTags(nil)
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, username string) (*model.User, error) {
	return auth.FetchUser(ctx, r.DB, username)
}

// AllocatedNodes is the resolver for the allocatedNodes field.
func (r *queryResolver) AllocatedNodes(ctx context.Context, cluster string) ([]*model.Count, error) {
	data, err := r.Repo.AllocatedNodes(cluster)
	if err != nil {
		return nil, err
	}

	counts := make([]*model.Count, 0, len(data))
	for subcluster, hosts := range data {
		counts = append(counts, &model.Count{
			Name:  subcluster,
			Count: len(hosts),
		})
	}

	return counts, nil
}

// Job is the resolver for the job field.
func (r *queryResolver) Job(ctx context.Context, id string) (*schema.Job, error) {
	numericId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	job, err := r.Repo.FindById(numericId)
	if err != nil {
		return nil, err
	}

	if user := auth.GetUser(ctx); user != nil && !user.HasRole(auth.RoleAdmin) && !user.HasRole(auth.RoleSupport) && job.User != user.Username {
		return nil, errors.New("you are not allowed to see this job")
	}

	return job, nil
}

// JobMetrics is the resolver for the jobMetrics field.
func (r *queryResolver) JobMetrics(ctx context.Context, id string, metrics []string, scopes []schema.MetricScope) ([]*model.JobMetricWithName, error) {
	job, err := r.Query().Job(ctx, id)
	if err != nil {
		return nil, err
	}

	data, err := metricdata.LoadData(job, metrics, scopes, ctx)
	if err != nil {
		return nil, err
	}

	res := []*model.JobMetricWithName{}
	for name, md := range data {
		for scope, metric := range md {
			res = append(res, &model.JobMetricWithName{
				Name:   name,
				Scope:  scope,
				Metric: metric,
			})
		}
	}

	return res, err
}

// JobsFootprints is the resolver for the jobsFootprints field.
func (r *queryResolver) JobsFootprints(ctx context.Context, filter []*model.JobFilter, metrics []string) (*model.Footprints, error) {
	return r.jobsFootprints(ctx, filter, metrics)
}

// Jobs is the resolver for the jobs field.
func (r *queryResolver) Jobs(ctx context.Context, filter []*model.JobFilter, page *model.PageRequest, order *model.OrderByInput) (*model.JobResultList, error) {
	if page == nil {
		page = &model.PageRequest{
			ItemsPerPage: 50,
			Page:         1,
		}
	}

	jobs, err := r.Repo.QueryJobs(ctx, filter, page, order)
	if err != nil {
		return nil, err
	}

	count, err := r.Repo.CountJobs(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &model.JobResultList{Items: jobs, Count: &count}, nil
}

// JobsStatistics is the resolver for the jobsStatistics field.
func (r *queryResolver) JobsStatistics(ctx context.Context, filter []*model.JobFilter, groupBy *model.Aggregate) ([]*model.JobsStatistics, error) {
	return r.jobsStatistics(ctx, filter, groupBy)
}

// JobsCount is the resolver for the jobsCount field.
func (r *queryResolver) JobsCount(ctx context.Context, filter []*model.JobFilter, groupBy model.Aggregate, weight *model.Weights, limit *int) ([]*model.Count, error) {
	counts, err := r.Repo.CountGroupedJobs(ctx, groupBy, filter, weight, limit)
	if err != nil {
		return nil, err
	}

	res := make([]*model.Count, 0, len(counts))
	for name, count := range counts {
		res = append(res, &model.Count{
			Name:  name,
			Count: count,
		})
	}
	return res, nil
}

// RooflineHeatmap is the resolver for the rooflineHeatmap field.
func (r *queryResolver) RooflineHeatmap(ctx context.Context, filter []*model.JobFilter, rows int, cols int, minX float64, minY float64, maxX float64, maxY float64) ([][]float64, error) {
	return r.rooflineHeatmap(ctx, filter, rows, cols, minX, minY, maxX, maxY)
}

// NodeMetrics is the resolver for the nodeMetrics field.
func (r *queryResolver) NodeMetrics(ctx context.Context, cluster string, nodes []string, scopes []schema.MetricScope, metrics []string, from time.Time, to time.Time) ([]*model.NodeMetrics, error) {
	user := auth.GetUser(ctx)
	if user != nil && !user.HasRole(auth.RoleAdmin) {
		return nil, errors.New("you need to be an administrator for this query")
	}

	if metrics == nil {
		for _, mc := range archive.GetCluster(cluster).MetricConfig {
			metrics = append(metrics, mc.Name)
		}
	}

	data, err := metricdata.LoadNodeData(cluster, metrics, nodes, scopes, from, to, ctx)
	if err != nil {
		return nil, err
	}

	nodeMetrics := make([]*model.NodeMetrics, 0, len(data))
	for hostname, metrics := range data {
		host := &model.NodeMetrics{
			Host:    hostname,
			Metrics: make([]*model.JobMetricWithName, 0, len(metrics)*len(scopes)),
		}
		host.SubCluster, _ = archive.GetSubClusterByNode(cluster, hostname)

		for metric, scopedMetrics := range metrics {
			for _, scopedMetric := range scopedMetrics {
				host.Metrics = append(host.Metrics, &model.JobMetricWithName{
					Name:   metric,
					Scope:  schema.MetricScopeNode, // NodeMetrics allow fixed scope?
					Metric: scopedMetric,
				})
			}
		}

		nodeMetrics = append(nodeMetrics, host)
	}

	return nodeMetrics, nil
}

// NumberOfNodes is the resolver for the numberOfNodes field.
func (r *subClusterResolver) NumberOfNodes(ctx context.Context, obj *schema.SubCluster) (int, error) {
	nodeList, err := archive.ParseNodeList(obj.Nodes)
	if err != nil {
		return 0, err
	}
	// log.Debugf(">>>> See raw list definition here: %v", nodeList)
	stringList := nodeList.PrintList()
	// log.Debugf(">>>> See parsed list here: %v", stringList)
	numOfNodes := len(stringList)
	// log.Debugf(">>>> See numOfNodes here: %v", len(stringList))
	return numOfNodes, nil
}

// Cluster returns generated.ClusterResolver implementation.
func (r *Resolver) Cluster() generated.ClusterResolver { return &clusterResolver{r} }

// Job returns generated.JobResolver implementation.
func (r *Resolver) Job() generated.JobResolver { return &jobResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// SubCluster returns generated.SubClusterResolver implementation.
func (r *Resolver) SubCluster() generated.SubClusterResolver { return &subClusterResolver{r} }

type clusterResolver struct{ *Resolver }
type jobResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subClusterResolver struct{ *Resolver }
