begin;

alter table clients drop column merge_request_review_timeout;
alter table clients drop column merge_request_reviewers_count;
alter table clients drop column merge_request_review_mention;

commit;
