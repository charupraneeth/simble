CREATE TABLE IF NOT EXISTS sites (
    id SERIAL PRIMARY KEY,
    domain varchar(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
) ;

create table if not exists analytics(
	id BIGSERIAL primary key,
	site_id integer references sites(id) on delete cascade not null,
	
	visitor_id varchar(16) not null,
	
	path text not null default '/',
	
	browser_name text,
	device_type varchar(255),
	os_name text,
	
	country_code char(2),
	city_name text,
	
	created_at timestamptz default NOW() not null
) ;

create index idx_analytics_site_date on analytics (site_id, created_at desc) ;

CREATE INDEX idx_analytics_visitor_site ON analytics (site_id, visitor_id, created_at) ;
