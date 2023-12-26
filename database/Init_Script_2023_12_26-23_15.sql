

-- Create Permissions Table
CREATE TABLE permission_types (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(50),
  description varchar(225),
  created_user bigint unsigned DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_user bigint unsigned DEFAULT NULL,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP,
  deleted_user bigint unsigned DEFAULT NULL,
  deleted_at datetime DEFAULT NULL,
  PRIMARY KEY (id)
);

-- Insert example data for Permissions Table
INSERT INTO permission_types (name, description, created_user, updated_at) VALUES
('view_user', 'Allow user to view other user and indivual user.', 1, NULL),
('create_user', 'Allow user to create other users.', 1, NULL),
('edit_user', 'Allow user to edit other user accounts and indivual user account.', 1, NULL),
('delete_user', 'Allow user to soft delete an indivual user account.', 1, NULL),
('view_role', 'Allow user to view system roles and indivual system role.', 1, NULL),
('create_role', 'Allow user to create system roles.', 1, NULL),
('edit_role', 'Allow user to edit system roles and indivual system role.', 1, NULL),
('delete_role', 'Allow user to soft delete system roles.', 1, NULL),
('view_product', 'Allow user to view product information', 1, NULL),
('create_product', 'Allow user to create product information', 1, NULL),
('edit_product', 'Allow user to edit product information', 1, NULL),
('delete_product', 'Allow user to delete products', 1, NULL),
('view_order', 'Allow user to view an order', 1, NULL),
('create_order', 'Allow user to create an order', 1, NULL),
('edit_order', 'Allow user to edit an order', 1, NULL),
('delete_order', 'Allow user to soft delete an order', 1, NULL),
('view_permission', 'Allow user to view permissions', 1, NULL),
('create_permission', 'Allow user to create permission', 1, NULL),
('edit_permission', 'Allow user to edit permission', 1, NULL),
('delete_permission', 'Allow user to delete permission', 1, NULL);

-- Create Roles Table
CREATE TABLE role_types (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(50),
  description varchar(225),
  created_user bigint unsigned DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_user bigint unsigned DEFAULT NULL,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP,
  deleted_user bigint unsigned DEFAULT NULL,
  deleted_at datetime DEFAULT NULL,
  PRIMARY KEY (id)
);

-- Insert example data for Roles Table
INSERT INTO role_types (name, description, created_user, updated_at) VALUES
('Admin', 'Should have all permissions i.e. [view_xxx, create_xxx, edit_xxx and delete_xxx] for [user, role, prodct, order and permission]', 1, NULL),
('Customer', 'Should have read and write permissions for orders i.e. [view_order, create_order, edit_order] and view permssions for products i.e. [view_product]', 1, NULL),
('Visitor', 'Should have just read permission for products i.e. [view_product]', 1, NULL);

-- Create Many-to-Many Table for Roles and Permissions
CREATE TABLE role_permissions (
  role_id bigint unsigned NOT NULL,
  permission_id bigint unsigned NOT NULL,
  created_user bigint unsigned DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_user bigint unsigned DEFAULT NULL,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP,
  deleted_user bigint unsigned DEFAULT NULL,
  deleted_at datetime DEFAULT NULL,
  PRIMARY KEY (role_id, permission_id),
  CONSTRAINT fk_role_permissions_permission 
    FOREIGN KEY (permission_id) 
    REFERENCES permission_types (id),
  CONSTRAINT fk_role_permissions_role 
    FOREIGN KEY (role_id) 
    REFERENCES role_types (id)
);
-- Insert example data for Role-Permissions Table
INSERT INTO role_permissions (role_id, permission_id, created_user, updated_at) VALUES
-- Admin Role
(1, 1, 1, NULL), -- view_user
(1, 2, 1, NULL), -- create_user
(1, 3, 1, NULL), -- edit_user
(1, 4, 1, NULL), -- delete_user
(1, 5, 1, NULL), -- view_role
(1, 6, 1, NULL), -- create_role
(1, 7, 1, NULL), -- edit_role
(1, 8, 1, NULL), -- delete_role
(1, 9, 1, NULL), -- view_product
(1, 10, 1, NULL), -- create_product
(1, 11, 1, NULL), -- edit_product
(1, 12, 1, NULL), -- delete_product
(1, 13, 1, NULL), -- view_order
(1, 14, 1, NULL), -- create_order
(1, 15, 1, NULL), -- edit_order
(1, 16, 1, NULL), -- delete_order
(1, 17, 1, NULL), -- view_permission
(1, 18, 1, NULL), -- create_permission
(1, 19, 1, NULL), -- edit_permission
(1, 20, 1, NULL), -- delete_permission

-- Customer Role
(2, 13, 1, NULL), -- view_order
(2, 14, 1, NULL), -- create_order
(2, 15, 1, NULL), -- edit_order
(2, 9, 1, NULL),  -- view_product

-- Visitor Role
(3, 9, 1, NULL);  -- view_product

-- Create Users Table
CREATE TABLE users (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  first_name varchar(50),
  last_name varchar(50),
  email varchar(225) DEFAULT NULL,
  hashed_password longblob,
  role_id bigint unsigned DEFAULT NULL,
  created_user bigint unsigned DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_user bigint unsigned DEFAULT NULL,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP,
  deleted_user bigint unsigned DEFAULT NULL,
  deleted_at datetime DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY email (email),
  CONSTRAINT fk_users_role 
  	FOREIGN KEY (role_id)
  	REFERENCES role_types (id)
);

-- Insert example data for Users Table
INSERT INTO users (first_name, last_name, email, hashed_password, role_id, created_user, updated_at) VALUES
('Admin', 'Doe', 'admin.doe@example.com', 'hashed_password', 1, NULL),
('Customer', 'Doe', 'customr.doe@example.com', 'hashed_password', 2, 1, NULL);

-- Create Products Table
CREATE TABLE products (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  title varchar(50),
  description varchar(225),
  price double DEFAULT NULL,
  created_user bigint unsigned DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_user bigint unsigned DEFAULT NULL,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP,
  deleted_user bigint unsigned DEFAULT NULL,
  deleted_at datetime DEFAULT NULL,
  PRIMARY KEY (id)
);

-- Insert example data for Products Table
INSERT INTO products (title, description, price, created_user, updated_at) VALUES
('Product 1', 'Description for Product 1', 29.99, 1, NULL),
('Product 2', 'Description for Product 2', 19.99, 2, NULL);

-- Create Pictures Table
CREATE TABLE pictures (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  picture_url varchar(50),
  product_id bigint unsigned NOT NULL,
  created_user bigint unsigned DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_user bigint unsigned DEFAULT NULL,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP,
  deleted_user bigint unsigned DEFAULT NULL,
  deleted_at datetime DEFAULT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_products_pictures 
  	FOREIGN KEY (product_id) 
  	REFERENCES products (id)
);

-- Insert example data for Pictures Table
INSERT INTO pictures (picture_url, product_id, created_user, updated_at) VALUES
('url1.jpg', 1, 1, NULL),
('url2.jpg', 2, 2 NULL);

-- Create Delivery Details Table
CREATE TABLE delivery_details (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  street_number varchar(50),
  street_name varchar(225),
  complex_name varchar(50),
  area_name varchar(50),
  city varchar(50),
  country varchar(50),
  desired_time datetime DEFAULT NULL,
  notes varchar(225),
  fullfilled_time datetime DEFAULT NULL,
  created_user bigint unsigned DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_user bigint unsigned DEFAULT NULL,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP,
  deleted_user bigint unsigned DEFAULT NULL,
  deleted_at datetime DEFAULT NULL,
  PRIMARY KEY (id)
);

-- Create Order Status Types Table
CREATE TABLE order_status_types (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(50),
  description varchar(225),
  created_user bigint unsigned DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_user bigint unsigned DEFAULT NULL,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP,
  deleted_user bigint unsigned DEFAULT NULL,
  deleted_at datetime DEFAULT NULL,
  PRIMARY KEY (id)
);

-- Insert example data for Order Status Types Table
INSERT INTO order_status_types (name, description, created_user, updated_at) VALUES
('Awaiting Payment', 'Order is awaiting payment', 1, NULL),
('Pending', 'Order is awaiting processing and confirmation at warehouses', 1, NULL),
('Shipped', 'Order has been shipped from warehouse', 1, NULL),
('Out for Delivery', 'Order out for delivery', 1, NULL),
('Order Complete', 'Order has completed all processes', 1, NULL);

-- Create Orders Table
CREATE TABLE orders (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  first_name varchar(50),
  last_name varchar(50),
  email varchar(225),
  status_id bigint unsigned DEFAULT NULL,
  delivery_details_id bigint unsigned NOT NULL,
  created_user bigint unsigned DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_user bigint unsigned DEFAULT NULL,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP,
  deleted_user bigint unsigned DEFAULT NULL,
  deleted_at datetime DEFAULT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_orders_delivery_details 
  	FOREIGN KEY (delivery_details_id) 
  	REFERENCES delivery_details (id),
  CONSTRAINT fk_orders_status_type 
  	FOREIGN KEY (status_id) 
  	REFERENCES order_status_types (id)
);

-- Create Order Items Table
CREATE TABLE order_items (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  order_id bigint unsigned DEFAULT NULL,
  product_title longtext,
  price float DEFAULT NULL,
  quantity bigint unsigned DEFAULT NULL,
  created_user bigint unsigned DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_user bigint unsigned DEFAULT NULL,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP,
  deleted_user bigint unsigned DEFAULT NULL,
  deleted_at datetime DEFAULT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_orders_order_items 
  	FOREIGN KEY (order_id) 
  	REFERENCES orders (id)
);