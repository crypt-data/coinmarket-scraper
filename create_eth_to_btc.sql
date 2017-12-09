create table if not exists eth_to_btc (

	time       integer not null,
	close      real    not null,
	high       real    not null,
	low        real    not null,
	open       real    not null,
	volumefrom real    not null,
	volumeto   real    not null,

primary key (time))
