begin;

alter table clients add column merge_request_old_timeout varchar(10) not null default '';
alter table clients add column merge_request_old_mention varchar(100) not null default '';

commit;
