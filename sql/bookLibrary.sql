
create database book_library;
use book_library;

CREATE TABLE IF NOT EXISTS BookCategories (
  id int(5) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  category varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Books (
  id int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
	isbn varchar(20) NOT NULL,
  title varchar(1000) NOT NULL,
  author  varchar(1000) NOT NULL,
  publication_year year NOT NULL,
  category_id int(5) NOT NULL,
  FOREIGN KEY (category_id) REFERENCES BookCategories(id)
);

CREATE TABLE IF NOT EXISTS Users (
  id int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  first_name varchar(512) NOT NULL,
  last_name varchar(512) NOT NULL,
  created_year year NOT NULL,
  phone_number int NOT NULL,
  login varchar(512) NOT NULL,
  password varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Employees (
  id int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  first_name varchar(512) NOT NULL,
  last_name varchar(512) NOT NULL,
  login varchar(512) NOT NULL,
  password varchar(255) NOT NULL,
  employee_type enum('librarian', 'admin') NOT NULL
);

CREATE TABLE IF NOT EXISTS BookCopies (
	id int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  book_id int(5) NOT NULL,
  status enum('avaliable', 'borrowed') DEFAULT 'avaliable',
  FOREIGN KEY (book_id) REFERENCES Books(id)
);

CREATE TABLE IF NOT EXISTS BorrowLog (
  id int(5) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  user_id int(11) NOT NULL,
  book_id int(11) NOT NULL,
  employee_id int(11) NOT NULL,
  release_date date NOT NULL,
  due_date date NOT NULL,
	returned_date date,
	returned boolean DEFAULT FALSE,
  FOREIGN KEY (user_id) REFERENCES Users(id),
  FOREIGN KEY (book_id) REFERENCES BookCopies(id),
  FOREIGN KEY (employee_id) REFERENCES Employees(id)
);

-- # Procedures

DELIMITER $$
CREATE PROCEDURE addCategory(
IN category varchar(255)
)
BEGIN

	DECLARE tempCategory varchar(255);
	DECLARE countCategory int;
	DECLARE successFlag int;

	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		ROLLBACK;
		SET @successFlag = 0;
		SELECT @successFlag;
	END;
	
	SET @tempCategory = category;
	SET @countCategory = (SELECT COUNT(*) FROM BookCategories WHERE category = tempCategory);
	
	SET autocommit = 0;
	START TRANSACTION;
		
		IF @countCategory = 0 THEN
		
			SET @query = 'INSERT INTO BookCategories(category) VALUES(?)';
			PREPARE stmt FROM @query;
			EXECUTE stmt USING @tempCategory;
			DEALLOCATE PREPARE stmt;
		END IF;
		
	COMMIT;

	SET @successFlag = 1;
	SELECT @successFlag;
	
END$$
DELIMITER ;

DELIMITER $$
CREATE PROCEDURE addBookCopy (
IN isbn varchar(20)
)
BEGIN 
	DECLARE tempISBN varchar(20);
	DECLARE bookId int(11);
	DECLARE successFlag int;

	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		ROLLBACK;
		SET @successFlag = 0;
		SELECT @successFlag;
	END;

	SET @tempISBN = isbn;
	SET @bookId = NULL;

	SET autocommit = 0;
	START TRANSACTION;
		SET @query = "SELECT id INTO @bookId FROM Books WHERE Books.isbn = ?";
		PREPARE stmt FROM @query;
		EXECUTE stmt USING @tempISBN;
		DEALLOCATE PREPARE stmt;

		SET @query2 = 'INSERT INTO BookCopies(book_id) VALUES(?)';
		PREPARE stmt FROM @query2;
		EXECUTE stmt USING @bookId;
		DEALLOCATE PREPARE stmt;
	COMMIT;

	SET @successFlag = 1;
	SELECT @successFlag;

END$$

DELIMITER $$
CREATE PROCEDURE addBook(
IN isbn varchar(20),
IN title varchar(1000),
IN author varchar(1000),
IN publication_year YEAR,
IN category varchar(255)
)
BEGIN
	DECLARE tempISBN varchar(20);
	DECLARE tempTitle VARCHAR(1000);
	DECLARE tempAuthor VARCHAR(1000);
	DECLARE tempPublicationYear YEAR;
	DECLARE tempCategory VARCHAR(255);
	DECLARE tempCategoryId int(5);
	DECLARE bookId int(11);
	DECLARE successFlag int;

	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		ROLLBACK;
		SET @successFlag = 0;
		SELECT @successFlag;
	END;

	SET @tempISBN = isbn;
	SET @tempTitle = title;
	SET @tempAuthor = author;
	SET @tempPublicationYear = publication_year;
	SET @tempCategory = category;
	SET @tempCategoryId = 0;
	
		
	SET autocommit = 0;
	START TRANSACTION;
	
		SET @queryCategory = 'SELECT id INTO @tempCategoryId FROM BookCategories WHERE BookCategories.category = ?';
		PREPARE stmt FROM @queryCategory;
		EXECUTE stmt USING @tempCategory;
		DEALLOCATE PREPARE stmt;
	
		SET @query1 = 'INSERT INTO Books(isbn, title, author, publication_year, category_id) VALUES(?, ?, ?, ?, ?)';
		PREPARE stmt FROM @query1;
		EXECUTE stmt USING @tempISBN, @tempTitle, @tempAuthor, @tempPublicationYear, @tempCategoryId;
		DEALLOCATE PREPARE stmt;

		SET @bookId = (SELECT lastBookID());

		SET @query2 = 'INSERT INTO BookCopies(book_id) VALUES(?)';
		PREPARE stmt FROM @query2;
		EXECUTE stmt USING @bookId;
		DEALLOCATE PREPARE stmt;
		
	COMMIT;
	
	SET @successFlag = 1;
	SELECT @successFlag;
	
END$$
DELIMITER ;


DELIMITER $$
CREATE PROCEDURE borrowBook(
IN user_id int(11),
IN book_id int(11),
IN employee_id int(11)
)
BEGIN
	DECLARE tempBookId int(11);
	DECLARE tempUserId int(11);
	DECLARE tempEmployeeId int(11);
	DECLARE checkIfEmployeeExists varchar(512);
	DECLARE checkIfUserExists varchar(512);
	DECLARE currDate date;
	DECLARE dueDate date;
	DECLARE successFlag INT;

	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		ROLLBACK;
		SET @successFlag = 0;
		SELECT @successFlag;
	END;
	
	SET @currDate = (SELECT CAST(CURRENT_TIMESTAMP AS DATE));
	SET @dueDate = (SELECT DATE_ADD(CAST(CURRENT_TIMESTAMP AS DATE), INTERVAL 1 MONTH));
	SET @tempBookId = book_id;
	SET @tempUserId = user_id;
	SET @tempEmployeeId = employee_id;
	
	SET autocommit = 0;	
	START TRANSACTION;

		SET @query = 'SELECT first_name INTO @checkIfEmployeeExists FROM Employees WHERE Employees.id = ?';
		PREPARE stmt FROM @query;
		EXECUTE stmt USING @tempEmployeeId;
		DEALLOCATE PREPARE stmt;

		IF @checkIfEmployeeExists IS NOT NULL THEN

			SET @query = 'SELECT first_name INTO @checkIfUserExists FROM Users WHERE Users.id = ?';
			PREPARE stmt FROM @query;
			EXECUTE stmt USING @tempUserId;
			DEALLOCATE PREPARE stmt;

			IF @checkIfUserExists IS NOT NULL THEN

				SET @query = 'SELECT status INTO @tempStatus FROM BookCopies WHERE BookCopies.book_id = ?';
				PREPARE stmt FROM @query;
				EXECUTE stmt USING @tempBookId;
				DEALLOCATE PREPARE stmt;

				IF STRCMP(@tempStatus, 'avaliable') = 0 THEN

					SET @query = 'INSERT INTO BorrowLog(user_id, book_id, employee_id, release_date, due_date) VALUES(?, ?, ?, ?, ?)';
					PREPARE stmt FROM @query;
					EXECUTE stmt USING @tempUserId, @tempBookId, @tempEmployeeId, @currDate, @dueDate;
					DEALLOCATE PREPARE stmt;

					SET @successFlag = 1;

				ELSE
					SET @successFlag = 0;
				END IF;

			ELSE 
				SET @successFlag = 0;
			END IF;

		ELSE 
			SET @successFlag = 0;
		END IF;

		SELECT @successFlag;

	COMMIT;
	
END $$
DELIMITER ;

DELIMITER $$
CREATE PROCEDURE registerUser (
IN first_name varchar(512),
IN last_name varchar(512),
IN phone_number int,
IN login varchar(512), 
IN password varchar(255)
)
BEGIN
	DECLARE tempFirstName varchar(512);
	DECLARE tempLastName varchar(512);
	DECLARE tempPhoneNumber int;
	DECLARE tempLogin varchar(512);
	DECLARE tempPassword varchar(255);
	DECLARE currYear year;
	DECLARE successFlag INT;
	DECLARE countLogin INT DEFAULT NULL;
	DECLARE newSHA2 VARCHAR(255);

	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		ROLLBACK;
		SET @successFlag = 0;
		SELECT @successFlag;
	END;

	SET @tempFirstName = first_name;
	SET @tempLastName = last_name;
	SET @tempPhoneNumber = phone_number;
	SET @tempPassword = password;
	SET @tempLogin = login;
	SET @countLogin = (SELECT COUNT(*) FROM Users WHERE Users.login = @tempLogin);
	SET @currYear = (SELECT YEAR(CAST(CURRENT_TIMESTAMP AS DATE)));
	
	IF @countLogin <= 0 THEN
	
		SET @newSHA2 = SHA2(@tempPassword, 0);
	
		SET @query = 'INSERT INTO Users(first_name, last_name, created_year, phone_number, login, password) VALUES(?,?,?,?,?,?)';
		PREPARE stmt FROM @query;
		EXECUTE stmt USING @tempFirstName, @tempLastName, @currYear, @tempPhoneNumber, @tempLogin, @newSHA2;
		DEALLOCATE PREPARE stmt;
	
		SET successFlag = 1;
	
	ELSE
		SET successFlag = 0;
	END IF;	
	
	SELECT successFlag;
	
END$$
DELIMITER ;

DELIMITER $$
CREATE PROCEDURE registerEmployee (
IN first_name varchar(512),
IN last_name varchar(512),
IN login varchar(512), 
IN password varchar(255),
IN typed_type varchar(255)

)
BEGIN
	DECLARE tempFirstName varchar(512);
	DECLARE tempLastName varchar(512);
	DECLARE tempLogin varchar(512);
	DECLARE tempPassword varchar(255);
	DECLARE successFlag INT;
	DECLARE countLogin INT DEFAULT NULL;
	DECLARE newSHA2 VARCHAR(255);
	DECLARE type varchar(20);

	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		ROLLBACK;
		SET @successFlag = 0;
		SELECT @successFlag;
	END;
	
	SET @tempFirstName = first_name;
	SET @tempLastName = last_name;
	SET @tempPassword = password;
	SET @tempLogin = login;
	SET @type = typed_type;
	SET @countLogin = (SELECT COUNT(*) FROM Employees WHERE Employees.login = @tempLogin);
	
	IF @countLogin <= 0 THEN
	
		SET @newSHA2 = SHA2(@tempPassword, 0);
	
		SET @query = 'INSERT INTO Employees(first_name, last_name, login, password, employee_type) VALUES(?,?,?,?,?)';
		PREPARE stmt FROM @query;
		EXECUTE stmt USING @tempFirstName, @tempLastName, @tempLogin, @newSHA2, @type;
		DEALLOCATE PREPARE stmt;
	
		SET successFlag = 1;
	
	ELSE
		SET successFlag = 0;
	END IF;	
	
	SELECT successFlag;
	
END$$
DELIMITER ;

DELIMITER $$
CREATE PROCEDURE deleteUser(
IN user_id int(11)
)
BEGIN 
	DECLARE tempUserId INT(11);
	DECLARE tempCountId INT DEFAULT 0;
	DECLARE successFlag int;

	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		ROLLBACK;
		SET @successFlag = 0;
		SELECT @successFlag;
	END;
	
	SET @tempUserId = user_id;
	
	SET @query = 'SELECT COUNT(*) INTO @tempCountId FROM Users WHERE Users.id = ?';
	PREPARE stmt FROM @query;
	EXECUTE stmt USING @tempUserId;
	DEALLOCATE PREPARE stmt;
	
	START TRANSACTION;
	
		IF @tempCountId = 0 THEN
			SET @successFlag = 0;
		ELSE
			set @query = 'DELETE FROM Users WHERE Users.id = ?';
			PREPARE stmt FROM @query;
			EXECUTE stmt USING @tempUserId;
			DEALLOCATE PREPARE stmt;
			
			SET @successFlag = 1;
		END IF;
		SELECT @successFlag;
	COMMIT;

END$$
DELIMITER ;

DELIMITER $$
CREATE PROCEDURE deleteEmployee(
IN employee_id int(11)
)
BEGIN 
	DECLARE tempEmployeeId INT(11);
	DECLARE tempCountId INT DEFAULT 0;
	DECLARE successFlag int;

	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		ROLLBACK;
		SET @successFlag = 0;
		SELECT @successFlag;
	END;
	
	SET @tempEmployeeId = employee_id;
	
	SET @query = 'SELECT COUNT(*) INTO @tempCountId FROM Employees WHERE Employees.id = ?';
	PREPARE stmt FROM @query;
	EXECUTE stmt USING @tempEmployeeId;
	DEALLOCATE PREPARE stmt;
	
	START TRANSACTION;
	
		IF @tempCountId <= 0 THEN
			SET @successFlag = 0;
		ELSE
			SET @query = 'DELETE FROM Employees WHERE Employees.id = ?';
			PREPARE stmt FROM @query;
			EXECUTE stmt USING @tempEmployeeId;
			DEALLOCATE PREPARE stmt;
			
			SET @successFlag = 1;
		END IF;
		SELECT @successFlag;
	COMMIT;

END$$
DELIMITER ;

DELIMITER $$
CREATE PROCEDURE returnBook(
IN user_id int(11),
IN book_id int(11)
)
BEGIN 
	DECLARE tempUserId int(11);
	DECLARE tempBookId int(11);
	DECLARE tempCountLogs INT DEFAULT 0;
	DECLARE returnDate date;
	DECLARE successFlag int;

	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		ROLLBACK;
		SET @successFlag = 0;
		SELECT @successFlag;
	END;
	
	SET @tempUserId = user_id;
	SET @tempBookId = book_id;
	SET @retunDate = (SELECT CAST(CURRENT_TIMESTAMP AS DATE));
	
	SET @query = 'SELECT COUNT(*) INTO @tempCountLogs FROM BorrowLog WHERE BorrowLog.user_id = ? AND BorrowLog.book_id = ?';
	PREPARE stmt FROM @query;
	EXECUTE stmt USING @tempUserId, @tempBookId;
	DEALLOCATE PREPARE stmt;
	
	START TRANSACTION;
	
		IF @tempCountLogs = 0 THEN
			SET @successFlag = 0;
			
		ELSE
			SET @query = 'UPDATE BorrowLog SET BorrowLog.returned_date = ? WHERE BorrowLog.user_id = ? AND BorrowLog.book_id = ?';
			PREPARE stmt FROM @query;
			EXECUTE stmt USING @retunDate, @tempUserId, @tempBookId;
			DEALLOCATE PREPARE stmt;

			SET @query = 'UPDATE BorrowLog SET BorrowLog.returned = TRUE WHERE BorrowLog.user_id = ? AND BorrowLog.book_id = ?';
			PREPARE stmt FROM @query;
			EXECUTE stmt USING @tempUserId, @tempBookId;
			DEALLOCATE PREPARE stmt;
			
			SET @successFlag = 1;
		END IF;
		
		SELECT @successFlag;
	COMMIT;

END$$
DELIMITER ;



DELIMITER $$
CREATE PROCEDURE logInUser(
IN typed_login VARCHAR(512),
IN typed_password VARCHAR(255)
)
BEGIN
	DECLARE tempSha2 VARCHAR(255);
	DECLARE tempLogin VARCHAR(512);
	DECLARE tempHashPasswd VARCHAR(255);
	DECLARE successFlag INT;

	SET @tempHashPasswd = '';
	SET @tempLogin = typed_login;
	SET @tempPassword = typed_password;

	SET @hashQuery = 'SELECT password INTO @tempHashPasswd FROM Users WHERE Users.login = ? ';
	PREPARE stmt FROM @hashQuery;
	EXECUTE stmt USING @tempLogin;
	DEALLOCATE PREPARE stmt;
	SET @tempSha2 = SHA2(@tempPassword, 0);

	IF STRCMP(@tempSha2, @tempHashPasswd) = 0 THEN
		SET @successFlag = 1;
	ELSE
		SET @successFlag = 0;
	END IF;

	SELECT @successFlag;
END$$
DELIMITER ;



DELIMITER $$
CREATE PROCEDURE logInLibrarian(
IN typed_login VARCHAR(512),
IN typed_password VARCHAR(255),
IN typed_type VARCHAR(255)
)
BEGIN 
	DECLARE tempSha2 VARCHAR(255); 
	DECLARE tempLogin VARCHAR(512); 
	DECLARE tempHashPwd VARCHAR(255);
	DECLARE successFlag INT; 
	DECLARE type varchar(20);  

	SET @tempHashPwd = ''; 
	SET @tempLogin = typed_login; 
	SET @tempPassword = typed_password; 
	SET @type = typed_type; 

	SET @hashQuery = 'SELECT password INTO @tempHashPwd FROM Employees WHERE Employees.login = ? AND employee_type = ?'; 
	PREPARE stmt FROM @hashQuery; 
	EXECUTE stmt USING @tempLogin, @type; 
	DEALLOCATE PREPARE stmt; 
	
	SET @tempSha2 = SHA2(@tempPassword, 0); 

	IF STRCMP(@tempSha2, @tempHashPwd) = 0 THEN 
		SET @successFlag = 1; 
	ELSE 
		SET @successFlag = 0; 
	END IF; 

	SELECT @successFlag; 
END$$
DELIMITER ;


DELIMITER $$
CREATE PROCEDURE renewBook (
	IN user_id int(11),
	IN book_id int(11)
)
BEGIN
	DECLARE tempUserId int(11);
	DECLARE tempBookId int(11);
	DECLARE releaseDate date;
	DECLARE dueDate date;
	DECLARE newDate date;
	DECLARE monthDiff int;
	DECLARE logId int;

	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		ROLLBACK;
		SET @successFlag = 0;
		SELECT @successFlag;
	END;

	SET @tempUserId = user_id;
	SET @tempBookId = book_id;
	SET @monthDiff = 0;
	SET @logId = 0;
	SET @newDate = NULL;
	SET @releaseDate = NULL;
	SET @dueDate = NULL;

	START TRANSACTION;
		SET @query = 'SELECT BorrowLog.id, BorrowLog.release_date, BorrowLog.due_date INTO @logId, @releaseDate, @dueDate FROM BorrowLog WHERE BorrowLog.user_id = ? AND BorrowLog.book_id = ? AND BorrowLog.returned = FALSE';
		PREPARE stmt FROM @query;
		EXECUTE stmt USING @tempUserId, @tempBookId;
		DEALLOCATE PREPARE stmt;

		SELECT @releaseDate;
		SELECT @dueDate;

		IF @logId = 0 THEN
			SET @successFlag = 0;
		ELSE
			SET @newDate = (SELECT DATE_ADD(@dueDate, INTERVAL 1 MONTH));
			SELECT @newDate;
			SET @monthDiff = (SELECT TIMESTAMPDIFF(MONTH, @releaseDate, @newDate));
			SELECT @monthDiff;

			IF @monthDiff > 2 THEN
				SET @successFlag = 0;
			ELSE 

				SET @query = 'UPDATE BorrowLog SET BorrowLog.due_date = ? WHERE user_id = ? AND book_id = ? AND BorrowLog.returned = FALSE';
				PREPARE stmt FROM @query;
				EXECUTE stmt USING @newDate, @tempUserId, @tempBookId;
				DEALLOCATE PREPARE stmt;
				
				SET @successFlag = 1;

			END IF;
		END IF;
		
		SELECT @successFlag;
	COMMIT;
END$$
DELIMITER ;



# Functions

DELIMITER $$
CREATE FUNCTION lastBookID()
RETURNS int(11)
NOT DETERMINISTIC
BEGIN
	DECLARE lastID int(11);
	
	SET lastID = (SELECT COUNT(*) FROM Books);
	
	IF lastID = 0 THEN 
		SET lastID = 1;
		RETURN lastID;
	END IF;
	
	SET lastID = (SELECT MAX(id) FROM Books);
	
	RETURN lastID;
END$$
DELIMITER ;


--# Triggers

DELIMITER $$
CREATE TRIGGER changeStatusBorrowed AFTER INSERT
ON BorrowLog
FOR EACH ROW 
BEGIN

	UPDATE BookCopies SET BookCopies.status = 'borrowed' WHERE BookCopies.book_id = NEW.book_id;
	
END $$
DELIMITER ;

DELIMITER $$
CREATE TRIGGER changeStatusReturned  AFTER UPDATE
ON BorrowLog
FOR EACH ROW
BEGIN
	IF NEW.returned = TRUE THEN
		UPDATE BookCopies SET BookCopies.status = 'avaliable' WHERE BookCopies.book_id = OLD.book_id;
	END IF;

END $$
DELIMITER ;


# SimpleInstructions

CALL registerUser('Anna', 'Kowalski', 784254125, 'anka99', 'ak123');
CALL registerUser('Wiktor', 'Matejko', 653586452, 'wikmat', 'wm123');

CALL registerEmployee('Adam', 'Nowak', 'anowak', 'an123', 'librarian');
CALL registerEmployee('Jan', 'Kowalski', 'jkowalski', 'jk123', 'librarian');
CALL registerEmployee('Admin', 'Admin', 'admin', 'admin123', 'admin');

CALL addCategory('Horror');
CALL addCategory('Romance');
CALL addCategory('Science Fiction');
CALL addCategory('Adventure');
CALL addCategory('Fairy tale');
CALL addCategory('Crime');
CALL addCategory('Fantasy');
CALL addCategory('Historical fiction');

CALL addBook('978-3-16-148410-0','A Tale Of Two Cities', 'Charles Dickens', 2018, 'Historical fiction');
CALL addBook('922-3-16-148410-0', 'The Little Prince', 'Antoine de Saint-Exupery', 2021, 'Adventure');
CALL addBook('989-3-16-148410-0','The Hobbit', 'JRR Tolkien', 2015, 'Fantasy');
CALL addBook('998-3-22-148610-0', 'Death On The Nile', 'Agatha Christie', 2012, 'Crime');

--# Views

CREATE VIEW userView AS (
	SELECT id, first_name, last_name, created_year, phone_number FROM Users
);

CREATE VIEW employeeView AS (
	SELECT id, first_name, last_name FROM Employees
);

--# PrepareStatements

PREPARE showBooksBasedOnStatus FROM 'SELECT Books.title FROM Books JOIN BookCopies ON BookCopies.book_id = Books.id WHERE BookCopies.status = ?';

PREPARE showUserInfoBasedOnLogin FROM 'SELECT first_name, last_name, created_year, phone_number FROM Users WHERE Users.login = ?';

PREPARE showUsersBorrowsBasedOnLogin FROM 'SELECT title, release_date, due_date FROM BorrowLog JOIN Books ON BorrowLog.book_id = Books.id JOIN Users ON Users.id = BorrowLog.user_id WHERE Users.login = ?';

PREPARE showBooksBasedOnTitle FROM 'SELECT title, author, publication_year, category FROM Books JOIN BookCategories WHERE Books.title = ?';

PREPARE showBooksBasedOnKeyWordTitle FROM 'SELECT title, author, publication_year, category FROM Books JOIN BookCategories WHERE Books.title LIKE \'%?%\'';

PREPARE showBooksBasedOnKeyWordAuthor FROM 'SELECT title, author, publication_year, category FROM Books JOIN BookCategories WHERE Books.author LIKE \'%?%\'';

PREPARE showBooksBasedOnAuthor FROM 'SELECT title, author, publication_year, category FROM Books JOIN BookCategories WHERE Books.author = ?';

PREPARE showBooksBasedOnCategory FROM 'SELECT first_name, last_name, created_year, phone_number FROM Users WHERE Users.login = ?';

PREPARE showAvaliableCategories FROM 'SELECT category FROM BookCategories JOIN Books ON Books.category_id = BookCategories.id JOIN BookCopies ON BookCopies.book_id = Books.id WHERE BookCopies.status = \'avaliable\' GROUP BY BookCategories.id HAVING COUNT(*) > 0';

--# UsersPriviligies

CREATE USER 'admin'@'%';
SET PASSWORD FOR 'admin'@'%' = PASSWORD('admin');

CREATE USER 'librarian'@'%';
SET PASSWORD FOR 'librarian'@'%' = PASSWORD('librarian');

CREATE USER 'user'@'%';
SET PASSWORD FOR 'user'@'%' = PASSWORD('user');

-- Admin
GRANT ALL ON `book_library`.* TO 'admin'@'%';

-- User
GRANT EXECUTE ON PROCEDURE `book_library`.`registerUser` TO 'user'@'%';
GRANT EXECUTE ON PROCEDURE book_library.logInUser TO 'user'@'%';

GRANT SELECT ON book_library.userView to 'user'@'%';
GRANT SELECT ON book_library.BookCategories TO 'user'@'%';
GRANT SELECT ON book_library.BookCopies TO 'user'@'%';
GRANT SELECT ON book_library.Books TO 'user'@'%';

-- Librarian
GRANT EXECUTE ON PROCEDURE book_library.addBook TO 'librarian'@'%';
GRANT EXECUTE ON PROCEDURE book_library.borrowBook TO 'librarian'@'%';
GRANT EXECUTE ON PROCEDURE book_library.returnBook TO 'librarian'@'%';
GRANT EXECUTE ON PROCEDURE book_library.logInLibrarian TO 'librarian'@'%';

GRANT UPDATE, DELETE, INSERT, SELECT ON book_library.BorrowLog TO 'librarian'@'%';
GRANT UPDATE, DELETE, INSERT, SELECT ON book_library.Books TO 'librarian'@'%';
GRANT UPDATE, DELETE, INSERT, SELECT ON book_library.BookCopies TO 'librarian'@'%';
GRANT UPDATE, DELETE, INSERT, SELECT ON book_library.BookCategories TO 'librarian'@'%';
GRANT SELECT ON book_library.employeeView to 'librarian'@'%';
GRANT SELECT ON book_library.userView to 'librarian'@'%';

flush privileges;


