
CREATE TABLE [Todo]
(
    [id] [bigint] PRIMARY KEY NOT NULL,
    [title] [nvarchar](64) NOT NULL,
	status varchar NULL,

);

CREATE SEQUENCE MyTodoID
    START WITH 1
    INCREMENT BY 1
    NO CACHE
    ;
GO
