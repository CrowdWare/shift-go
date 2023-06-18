# Todo

- Releases should be build on the server using the same keys as the webservice
- The version of the lib should be packed in the shift.db first, last bytes or timestamp of the db file, then we can deploy multiple versions of the lib in the same apk and on the client the logic can choose which lib to use, to be backward compatible. First read version from dbfile then choose