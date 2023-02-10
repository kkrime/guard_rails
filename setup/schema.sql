CREATE DATABASE guard_rails;

\connect guard_rails

CREATE TABLE public.repositories(
    id SERIAL NOT NULL PRIMARY KEY,
	name CHARACTER VARYING(100) NOT NULL,
    url TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	deleted_at TIMESTAMP DEFAULT NULL
);
CREATE INDEX repositories_name ON public.repositories(name);
CREATE UNIQUE INDEX no_duplicate_repository_name ON public.repositories(name,deleted_at)
   WHERE deleted_at IS null;
CREATE UNIQUE INDEX no_duplicate_repository_urls  ON public.repositories(url)
   WHERE deleted_at IS null;

CREATE TYPE public.status AS ENUM (
	'QUEUED',
    'IN PROGRESS',
    'SUCCESS',
    'FAILURE'
);

CREATE TABLE public.scans(
    id SERIAL NOT NULL PRIMARY KEY,
    repository_id INTEGER REFERENCES public.repositories NOT NULL,
    status public.status DEFAULT 'QUEUED' NOT NULL,
    findings JSON,
	created_at TIMESTAMP DEFAULT NOW(),
	started_at TIMESTAMP DEFAULT NULL,
	ended_at   TIMESTAMP DEFAULT NULL
);
CREATE UNIQUE INDEX no_duplicate_queued_scans  ON public.scans(repository_id,status)
   WHERE status = 'QUEUED';
CREATE UNIQUE INDEX no_duplicate_in_progrss_scans  ON public.scans(repository_id,status)
   WHERE status = 'IN PROGRESS';
