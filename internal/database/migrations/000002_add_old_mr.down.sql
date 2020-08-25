begin;

alter table clients drop column merge_request_old_timeout;
alter table clients drop column merge_request_old_mention;

commit;
