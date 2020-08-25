begin;

alter table clients add column merge_request_review_timeout varchar(10) not null default '';
alter table clients add column merge_request_reviewers_count integer not null default 1;
alter table clients add column merge_request_review_mention varchar(100) not null default '';

commit;
