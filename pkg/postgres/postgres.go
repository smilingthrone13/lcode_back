package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func NewDBManger(db *pgxpool.Pool) *DbManager {
	return &DbManager{
		db: db,
	}
}

type DbManager struct {
	db *pgxpool.Pool
}

type TxKey struct{}

func (d *DbManager) TxOrDB(ctx context.Context) QueryExecutor {
	if tx, ok := ctx.Value(TxKey{}).(*Tx); ok {
		return tx.tx
	}

	return d.db
}

func (d *DbManager) GetDb() *pgxpool.Pool {
	return d.db
}

func NewTransactionProvider(db *pgxpool.Pool) *TransactionProvider {
	t := &TransactionProvider{db: db, txQueue: make(chan *Tx, 2000)}

	go t.afterSuccessWorker()

	return t
}

type TransactionProvider struct {
	db      *pgxpool.Pool
	txQueue chan *Tx
}

func (tm *TransactionProvider) afterSuccessWorker() {
	for tx := range tm.txQueue {
		for i := 0; i < len(tx.afterSuccessQueue); i++ {
			tx.afterSuccessQueue[i]()
		}
	}
}

func (tm *TransactionProvider) GetTxForParticipant(ctx context.Context) (TxForParticipant, error) {
	if tx, ok := ctx.Value(TxKey{}).(*Tx); ok {
		return tx, nil
	}

	return nil, errors.New("transaction not found in ctx")
}

func (tm *TransactionProvider) NewTx(ctx context.Context, opts *pgx.TxOptions) (Transaction, error) {
	var tx pgx.Tx
	var err error

	if opts == nil {
		tx, err = tm.db.Begin(ctx)
	} else {
		tx, err = tm.db.BeginTx(ctx, *opts)
	}

	if err != nil {
		err = errors.Wrap(err, "NewTx postgres pkg")

		return nil, err
	}

	return &Tx{tx: tx, txQueue: tm.txQueue}, nil
}

type Tx struct {
	tx                pgx.Tx
	afterSuccessQueue []func()
	txQueue           chan<- *Tx
}

func (t *Tx) Commit(ctx context.Context) error {
	if err := t.tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "Commit postgres pkg")
	}

	t.txQueue <- t

	return nil
}

func (t *Tx) AfterSuccess(ctx context.Context, f func()) {
	t.afterSuccessQueue = append(t.afterSuccessQueue, f)
}

func (t *Tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	AfterSuccess(ctx context.Context, f func())
}

type TxForParticipant interface {
	AfterSuccess(ctx context.Context, f func())
}

type QueryExecutor interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)

	//Get(dest interface{}, query string, args ...interface{}) error
	//GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	//
	//Select(dest interface{}, query string, args ...interface{}) error
	//SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	//
	//Query(query string, args ...any) (*sql.Rows, error)
	//QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	//QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	//
	//Exec(query string, args ...any) (sql.Result, error)
	//ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type errConstructor func(err error) error

type ErrPair struct {
	code           string
	errConstructor errConstructor
}

func Error(code string, err errConstructor) ErrPair {
	return ErrPair{code: code, errConstructor: err}
}

//func GetError(err error, errPairs ...ErrPair) error {
//	pqError, ok := err.(*pq.Error)
//	if !ok {
//		return err
//	}
//
//	for _, p := range errPairs {
//		if pqError.Code == pq.ErrorCode(p.code) {
//			return p.errConstructor(err)
//		}
//	}
//
//	return err
//}

//nolint:lll
const (

	// Section: Class 00 - Successful Completion.

	ERRCODE_SUCCESSFUL_COMPLETION = "00000" // successful_completion

	// Section: Class 01 - Warning.

	ERRCODE_WARNING                                       = "01000" // warning
	ERRCODE_WARNING_DYNAMIC_RESULT_SETS_RETURNED          = "0100C" // dynamic_result_sets_returned
	ERRCODE_WARNING_IMPLICIT_ZERO_BIT_PADDING             = "01008" // implicit_zero_bit_padding
	ERRCODE_WARNING_NULL_VALUE_ELIMINATED_IN_SET_FUNCTION = "01003" // null_value_eliminated_in_set_function
	ERRCODE_WARNING_PRIVILEGE_NOT_GRANTED                 = "01007" // privilege_not_granted
	ERRCODE_WARNING_PRIVILEGE_NOT_REVOKED                 = "01006" // privilege_not_revoked
	ERRCODE_WARNING_STRING_DATA_RIGHT_TRUNCATION          = "01004" // string_data_right_truncation
	ERRCODE_WARNING_DEPRECATED_FEATURE                    = "01P01" // deprecated_feature

	// Section: Class 02 - No Data (this is also a warning class per the SQL standard).

	ERRCODE_NO_DATA                                    = "02000" // no_data
	ERRCODE_NO_ADDITIONAL_DYNAMIC_RESULT_SETS_RETURNED = "02001" // no_additional_dynamic_result_sets_returned

	// Section: Class 03 - SQL Statement Not Yet Complete.

	ERRCODE_SQL_STATEMENT_NOT_YET_COMPLETE = "03000" // sql_statement_not_yet_complete

	// Section: Class 08 - Connection Exception.

	ERRCODE_CONNECTION_EXCEPTION                              = "08000" // connection_exception
	ERRCODE_CONNECTION_DOES_NOT_EXIST                         = "08003" // connection_does_not_exist
	ERRCODE_CONNECTION_FAILURE                                = "08006" // connection_failure
	ERRCODE_SQLCLIENT_UNABLE_TO_ESTABLISH_SQLCONNECTION       = "08001" // sqlclient_unable_to_establish_sqlconnection
	ERRCODE_SQLSERVER_REJECTED_ESTABLISHMENT_OF_SQLCONNECTION = "08004" // sqlserver_rejected_establishment_of_sqlconnection
	ERRCODE_TRANSACTION_RESOLUTION_UNKNOWN                    = "08007" // transaction_resolution_unknown
	ERRCODE_PROTOCOL_VIOLATION                                = "08P01" // protocol_violation

	// Section: Class 09 - Triggered Action Exception.

	ERRCODE_TRIGGERED_ACTION_EXCEPTION = "09000" // triggered_action_exception

	// Section: Class 0A - Feature Not Supported.

	ERRCODE_FEATURE_NOT_SUPPORTED = "0A000" // feature_not_supported

	// Section: Class 0B - Invalid Transaction Initiation.

	ERRCODE_INVALID_TRANSACTION_INITIATION = "0B000" // invalid_transaction_initiation

	// Section: Class 0F - Locator Exception.

	ERRCODE_LOCATOR_EXCEPTION         = "0F000" // locator_exception
	ERRCODE_L_E_INVALID_SPECIFICATION = "0F001" // invalid_locator_specification

	// Section: Class 0L - Invalid Grantor.

	ERRCODE_INVALID_GRANTOR         = "0L000" // invalid_grantor
	ERRCODE_INVALID_GRANT_OPERATION = "0LP01" // invalid_grant_operation

	// Section: Class 0P - Invalid Role Specification.

	ERRCODE_INVALID_ROLE_SPECIFICATION = "0P000" // invalid_role_specification

	// Section: Class 0Z - Diagnostics Exception.

	ERRCODE_DIAGNOSTICS_EXCEPTION                               = "0Z000" // diagnostics_exception
	ERRCODE_STACKED_DIAGNOSTICS_ACCESSED_WITHOUT_ACTIVE_HANDLER = "0Z002" // stacked_diagnostics_accessed_without_active_handler

	// Section: Class 20 - Case Not Found.

	ERRCODE_CASE_NOT_FOUND = "20000" // case_not_found

	// Section: Class 21 - Cardinality Violation.

	ERRCODE_CARDINALITY_VIOLATION = "21000" // cardinality_violation

	// Section: Class 22 - Data Exception.

	ERRCODE_DATA_EXCEPTION                                  = "22000" // data_exception
	ERRCODE_ARRAY_ELEMENT_ERROR                             = "2202E" // # SQL99's actual definition of "array element error" is subscript error
	ERRCODE_ARRAY_SUBSCRIPT_ERROR                           = "2202E" // array_subscript_error
	ERRCODE_CHARACTER_NOT_IN_REPERTOIRE                     = "22021" // character_not_in_repertoire
	ERRCODE_DATETIME_FIELD_OVERFLOW                         = "22008" // datetime_field_overflow
	ERRCODE_DATETIME_VALUE_OUT_OF_RANGE                     = "22008"
	ERRCODE_DIVISION_BY_ZERO                                = "22012" // division_by_zero
	ERRCODE_ERROR_IN_ASSIGNMENT                             = "22005" // error_in_assignment
	ERRCODE_ESCAPE_CHARACTER_CONFLICT                       = "2200B" // escape_character_conflict
	ERRCODE_INDICATOR_OVERFLOW                              = "22022" // indicator_overflow
	ERRCODE_INTERVAL_FIELD_OVERFLOW                         = "22015" // interval_field_overflow
	ERRCODE_INVALID_ARGUMENT_FOR_LOG                        = "2201E" // invalid_argument_for_logarithm
	ERRCODE_INVALID_ARGUMENT_FOR_NTILE                      = "22014" // invalid_argument_for_ntile_function
	ERRCODE_INVALID_ARGUMENT_FOR_NTH_VALUE                  = "22016" // invalid_argument_for_nth_value_function
	ERRCODE_INVALID_ARGUMENT_FOR_POWER_FUNCTION             = "2201F" // invalid_argument_for_power_function
	ERRCODE_INVALID_ARGUMENT_FOR_WIDTH_BUCKET_FUNCTION      = "2201G" // invalid_argument_for_width_bucket_function
	ERRCODE_INVALID_CHARACTER_VALUE_FOR_CAST                = "22018" // invalid_character_value_for_cast
	ERRCODE_INVALID_DATETIME_FORMAT                         = "22007" // invalid_datetime_format
	ERRCODE_INVALID_ESCAPE_CHARACTER                        = "22019" // invalid_escape_character
	ERRCODE_INVALID_ESCAPE_OCTET                            = "2200D" // invalid_escape_octet
	ERRCODE_INVALID_ESCAPE_SEQUENCE                         = "22025" // invalid_escape_sequence
	ERRCODE_NONSTANDARD_USE_OF_ESCAPE_CHARACTER             = "22P06" // nonstandard_use_of_escape_character
	ERRCODE_INVALID_INDICATOR_PARAMETER_VALUE               = "22010" // invalid_indicator_parameter_value
	ERRCODE_INVALID_PARAMETER_VALUE                         = "22023" // invalid_parameter_value
	ERRCODE_INVALID_PRECEDING_OR_FOLLOWING_SIZE             = "22013" // invalid_preceding_or_following_size
	ERRCODE_INVALID_REGULAR_EXPRESSION                      = "2201B" // invalid_regular_expression
	ERRCODE_INVALID_ROW_COUNT_IN_LIMIT_CLAUSE               = "2201W" // invalid_row_count_in_limit_clause
	ERRCODE_INVALID_ROW_COUNT_IN_RESULT_OFFSET_CLAUSE       = "2201X" // invalid_row_count_in_result_offset_clause
	ERRCODE_INVALID_TABLESAMPLE_ARGUMENT                    = "2202H" // invalid_tablesample_argument
	ERRCODE_INVALID_TABLESAMPLE_REPEAT                      = "2202G" // invalid_tablesample_repeat
	ERRCODE_INVALID_TIME_ZONE_DISPLACEMENT_VALUE            = "22009" // invalid_time_zone_displacement_value
	ERRCODE_INVALID_USE_OF_ESCAPE_CHARACTER                 = "2200C" // invalid_use_of_escape_character
	ERRCODE_MOST_SPECIFIC_TYPE_MISMATCH                     = "2200G" // most_specific_type_mismatch
	ERRCODE_NULL_VALUE_NOT_ALLOWED                          = "22004" // null_value_not_allowed
	ERRCODE_NULL_VALUE_NO_INDICATOR_PARAMETER               = "22002" // null_value_no_indicator_parameter
	ERRCODE_NUMERIC_VALUE_OUT_OF_RANGE                      = "22003" // numeric_value_out_of_range
	ERRCODE_SEQUENCE_GENERATOR_LIMIT_EXCEEDED               = "2200H" // sequence_generator_limit_exceeded
	ERRCODE_STRING_DATA_LENGTH_MISMATCH                     = "22026" // string_data_length_mismatch
	ERRCODE_STRING_DATA_RIGHT_TRUNCATION                    = "22001" // string_data_right_truncation
	ERRCODE_SUBSTRING_ERROR                                 = "22011" // substring_error
	ERRCODE_TRIM_ERROR                                      = "22027" // trim_error
	ERRCODE_UNTERMINATED_C_STRING                           = "22024" // unterminated_c_string
	ERRCODE_ZERO_LENGTH_CHARACTER_STRING                    = "2200F" // zero_length_character_string
	ERRCODE_FLOATING_POINT_EXCEPTION                        = "22P01" // floating_point_exception
	ERRCODE_INVALID_TEXT_REPRESENTATION                     = "22P02" // invalid_text_representation
	ERRCODE_INVALID_BINARY_REPRESENTATION                   = "22P03" // invalid_binary_representation
	ERRCODE_BAD_COPY_FILE_FORMAT                            = "22P04" // bad_copy_file_format
	ERRCODE_UNTRANSLATABLE_CHARACTER                        = "22P05" // untranslatable_character
	ERRCODE_NOT_AN_XML_DOCUMENT                             = "2200L" // not_an_xml_document
	ERRCODE_INVALID_XML_DOCUMENT                            = "2200M" // invalid_xml_document
	ERRCODE_INVALID_XML_CONTENT                             = "2200N" // invalid_xml_content
	ERRCODE_INVALID_XML_COMMENT                             = "2200S" // invalid_xml_comment
	ERRCODE_INVALID_XML_PROCESSING_INSTRUCTION              = "2200T" // invalid_xml_processing_instruction
	ERRCODE_DUPLICATE_JSON_OBJECT_KEY_VALUE                 = "22030" // duplicate_json_object_key_value
	ERRCODE_INVALID_ARGUMENT_FOR_SQL_JSON_DATETIME_FUNCTION = "22031" // invalid_argument_for_sql_json_datetime_function
	ERRCODE_INVALID_JSON_TEXT                               = "22032" // invalid_json_text
	ERRCODE_INVALID_SQL_JSON_SUBSCRIPT                      = "22033" // invalid_sql_json_subscript
	ERRCODE_MORE_THAN_ONE_SQL_JSON_ITEM                     = "22034" // more_than_one_sql_json_item
	ERRCODE_NO_SQL_JSON_ITEM                                = "22035" // no_sql_json_item
	ERRCODE_NON_NUMERIC_SQL_JSON_ITEM                       = "22036" // non_numeric_sql_json_item
	ERRCODE_NON_UNIQUE_KEYS_IN_A_JSON_OBJECT                = "22037" // non_unique_keys_in_a_json_object
	ERRCODE_SINGLETON_SQL_JSON_ITEM_REQUIRED                = "22038" // singleton_sql_json_item_required
	ERRCODE_SQL_JSON_ARRAY_NOT_FOUND                        = "22039" // sql_json_array_not_found
	ERRCODE_SQL_JSON_MEMBER_NOT_FOUND                       = "2203A" // sql_json_member_not_found
	ERRCODE_SQL_JSON_NUMBER_NOT_FOUND                       = "2203B" // sql_json_number_not_found
	ERRCODE_SQL_JSON_OBJECT_NOT_FOUND                       = "2203C" // sql_json_object_not_found
	ERRCODE_TOO_MANY_JSON_ARRAY_ELEMENTS                    = "2203D" // too_many_json_array_elements
	ERRCODE_TOO_MANY_JSON_OBJECT_MEMBERS                    = "2203E" // too_many_json_object_members
	ERRCODE_SQL_JSON_SCALAR_REQUIRED                        = "2203F" // sql_json_scalar_required

	// Section: Class 23 - Integrity Constraint Violation.

	ERRCODE_INTEGRITY_CONSTRAINT_VIOLATION = "23000" // integrity_constraint_violation
	ERRCODE_RESTRICT_VIOLATION             = "23001" // restrict_violation
	ERRCODE_NOT_NULL_VIOLATION             = "23502" // not_null_violation
	ERRCODE_FOREIGN_KEY_VIOLATION          = "23503" // foreign_key_violation
	ERRCODE_UNIQUE_VIOLATION               = "23505" // unique_violation
	ERRCODE_CHECK_VIOLATION                = "23514" // check_violation
	ERRCODE_EXCLUSION_VIOLATION            = "23P01" // exclusion_violation

	// Section: Class 24 - Invalid Cursor State.

	ERRCODE_INVALID_CURSOR_STATE = "24000" // invalid_cursor_state

	// Section: Class 25 - Invalid Transaction State.

	ERRCODE_INVALID_TRANSACTION_STATE                            = "25000" // invalid_transaction_state
	ERRCODE_ACTIVE_SQL_TRANSACTION                               = "25001" // active_sql_transaction
	ERRCODE_BRANCH_TRANSACTION_ALREADY_ACTIVE                    = "25002" // branch_transaction_already_active
	ERRCODE_HELD_CURSOR_REQUIRES_SAME_ISOLATION_LEVEL            = "25008" // held_cursor_requires_same_isolation_level
	ERRCODE_INAPPROPRIATE_ACCESS_MODE_FOR_BRANCH_TRANSACTION     = "25003" // inappropriate_access_mode_for_branch_transaction
	ERRCODE_INAPPROPRIATE_ISOLATION_LEVEL_FOR_BRANCH_TRANSACTION = "25004" // inappropriate_isolation_level_for_branch_transaction
	ERRCODE_NO_ACTIVE_SQL_TRANSACTION_FOR_BRANCH_TRANSACTION     = "25005" // no_active_sql_transaction_for_branch_transaction
	ERRCODE_READ_ONLY_SQL_TRANSACTION                            = "25006" // read_only_sql_transaction
	ERRCODE_SCHEMA_AND_DATA_STATEMENT_MIXING_NOT_SUPPORTED       = "25007" // schema_and_data_statement_mixing_not_supported
	ERRCODE_NO_ACTIVE_SQL_TRANSACTION                            = "25P01" // no_active_sql_transaction
	ERRCODE_IN_FAILED_SQL_TRANSACTION                            = "25P02" // in_failed_sql_transaction
	ERRCODE_IDLE_IN_TRANSACTION_SESSION_TIMEOUT                  = "25P03" // idle_in_transaction_session_timeout

	// Section: Class 26 - Invalid SQL Statement Name.

	ERRCODE_INVALID_SQL_STATEMENT_NAME = "26000" // invalid_sql_statement_name

	// Section: Class 27 - Triggered Data Change Violation.

	ERRCODE_TRIGGERED_DATA_CHANGE_VIOLATION = "27000" // triggered_data_change_violation

	// Section: Class 28 - Invalid Authorization Specification.

	ERRCODE_INVALID_AUTHORIZATION_SPECIFICATION = "28000" // invalid_authorization_specification
	ERRCODE_INVALID_PASSWORD                    = "28P01" // invalid_password

	// Section: Class 2B - Dependent Privilege Descriptors Still Exist.

	ERRCODE_DEPENDENT_PRIVILEGE_DESCRIPTORS_STILL_EXIST = "2B000" // dependent_privilege_descriptors_still_exist
	ERRCODE_DEPENDENT_OBJECTS_STILL_EXIST               = "2BP01" // dependent_objects_still_exist

	// Section: Class 2D - Invalid Transaction Termination.

	ERRCODE_INVALID_TRANSACTION_TERMINATION = "2D000" // invalid_transaction_termination

	// Section: Class 2F - SQL Routine Exception.

	ERRCODE_SQL_ROUTINE_EXCEPTION                       = "2F000" // sql_routine_exception
	ERRCODE_S_R_E_FUNCTION_EXECUTED_NO_RETURN_STATEMENT = "2F005" // function_executed_no_return_statement
	ERRCODE_S_R_E_MODIFYING_SQL_DATA_NOT_PERMITTED      = "2F002" // modifying_sql_data_not_permitted
	ERRCODE_S_R_E_PROHIBITED_SQL_STATEMENT_ATTEMPTED    = "2F003" // prohibited_sql_statement_attempted
	ERRCODE_S_R_E_READING_SQL_DATA_NOT_PERMITTED        = "2F004" // reading_sql_data_not_permitted

	// Section: Class 34 - Invalid Cursor Name.

	ERRCODE_INVALID_CURSOR_NAME = "34000" // invalid_cursor_name

	// Section: Class 38 - External Routine Exception.

	ERRCODE_EXTERNAL_ROUTINE_EXCEPTION               = "38000" // external_routine_exception
	ERRCODE_E_R_E_CONTAINING_SQL_NOT_PERMITTED       = "38001" // containing_sql_not_permitted
	ERRCODE_E_R_E_MODIFYING_SQL_DATA_NOT_PERMITTED   = "38002" // modifying_sql_data_not_permitted
	ERRCODE_E_R_E_PROHIBITED_SQL_STATEMENT_ATTEMPTED = "38003" // prohibited_sql_statement_attempted
	ERRCODE_E_R_E_READING_SQL_DATA_NOT_PERMITTED     = "38004" // reading_sql_data_not_permitted

	// Section: Class 39 - External Routine Invocation Exception.

	ERRCODE_EXTERNAL_ROUTINE_INVOCATION_EXCEPTION   = "39000" // external_routine_invocation_exception
	ERRCODE_E_R_I_E_INVALID_SQLSTATE_RETURNED       = "39001" // invalid_sqlstate_returned
	ERRCODE_E_R_I_E_NULL_VALUE_NOT_ALLOWED          = "39004" // null_value_not_allowed
	ERRCODE_E_R_I_E_TRIGGER_PROTOCOL_VIOLATED       = "39P01" // trigger_protocol_violated
	ERRCODE_E_R_I_E_SRF_PROTOCOL_VIOLATED           = "39P02" // srf_protocol_violated
	ERRCODE_E_R_I_E_EVENT_TRIGGER_PROTOCOL_VIOLATED = "39P03" // event_trigger_protocol_violated

	// Section: Class 3B - Savepoint Exception.

	ERRCODE_SAVEPOINT_EXCEPTION       = "3B000" // savepoint_exception
	ERRCODE_S_E_INVALID_SPECIFICATION = "3B001" // invalid_savepoint_specification

	// Section: Class 3D - Invalid Catalog Name.

	ERRCODE_INVALID_CATALOG_NAME = "3D000" // invalid_catalog_name

	// Section: Class 3F - Invalid Schema Name.

	ERRCODE_INVALID_SCHEMA_NAME = "3F000" // invalid_schema_name

	// Section: Class 40 - Transaction Rollback.

	ERRCODE_TRANSACTION_ROLLBACK               = "40000" // transaction_rollback
	ERRCODE_T_R_INTEGRITY_CONSTRAINT_VIOLATION = "40002" // transaction_integrity_constraint_violation
	ERRCODE_T_R_SERIALIZATION_FAILURE          = "40001" // serialization_failure
	ERRCODE_T_R_STATEMENT_COMPLETION_UNKNOWN   = "40003" // statement_completion_unknown
	ERRCODE_T_R_DEADLOCK_DETECTED              = "40P01" // deadlock_detected

	// Section: Class 42 - Syntax Error or Access Rule Violation.

	ERRCODE_SYNTAX_ERROR_OR_ACCESS_RULE_VIOLATION = "42000" // syntax_error_or_access_rule_violation
	ERRCODE_SYNTAX_ERROR                          = "42601" // syntax_error
	ERRCODE_INSUFFICIENT_PRIVILEGE                = "42501" // insufficient_privilege
	ERRCODE_CANNOT_COERCE                         = "42846" // cannot_coerce
	ERRCODE_GROUPING_ERROR                        = "42803" // grouping_error
	ERRCODE_WINDOWING_ERROR                       = "42P20" // windowing_error
	ERRCODE_INVALID_RECURSION                     = "42P19" // invalid_recursion
	ERRCODE_INVALID_FOREIGN_KEY                   = "42830" // invalid_foreign_key
	ERRCODE_INVALID_NAME                          = "42602" // invalid_name
	ERRCODE_NAME_TOO_LONG                         = "42622" // name_too_long
	ERRCODE_RESERVED_NAME                         = "42939" // reserved_name
	ERRCODE_DATATYPE_MISMATCH                     = "42804" // datatype_mismatch
	ERRCODE_INDETERMINATE_DATATYPE                = "42P18" // indeterminate_datatype
	ERRCODE_COLLATION_MISMATCH                    = "42P21" // collation_mismatch
	ERRCODE_INDETERMINATE_COLLATION               = "42P22" // indeterminate_collation
	ERRCODE_WRONG_OBJECT_TYPE                     = "42809" // wrong_object_type
	ERRCODE_GENERATED_ALWAYS                      = "428C9" // generated_always

	ERRCODE_UNDEFINED_COLUMN              = "42703" // undefined_column
	ERRCODE_UNDEFINED_CURSOR              = "34000"
	ERRCODE_UNDEFINED_DATABASE            = "3D000"
	ERRCODE_UNDEFINED_FUNCTION            = "42883" // undefined_function
	ERRCODE_UNDEFINED_PSTATEMENT          = "26000"
	ERRCODE_UNDEFINED_SCHEMA              = "3F000"
	ERRCODE_UNDEFINED_TABLE               = "42P01" // undefined_table
	ERRCODE_UNDEFINED_PARAMETER           = "42P02" // undefined_parameter
	ERRCODE_UNDEFINED_OBJECT              = "42704" // undefined_object
	ERRCODE_DUPLICATE_COLUMN              = "42701" // duplicate_column
	ERRCODE_DUPLICATE_CURSOR              = "42P03" // duplicate_cursor
	ERRCODE_DUPLICATE_DATABASE            = "42P04" // duplicate_database
	ERRCODE_DUPLICATE_FUNCTION            = "42723" // duplicate_function
	ERRCODE_DUPLICATE_PSTATEMENT          = "42P05" // duplicate_prepared_statement
	ERRCODE_DUPLICATE_SCHEMA              = "42P06" // duplicate_schema
	ERRCODE_DUPLICATE_TABLE               = "42P07" // duplicate_table
	ERRCODE_DUPLICATE_ALIAS               = "42712" // duplicate_alias
	ERRCODE_DUPLICATE_OBJECT              = "42710" // duplicate_object
	ERRCODE_AMBIGUOUS_COLUMN              = "42702" // ambiguous_column
	ERRCODE_AMBIGUOUS_FUNCTION            = "42725" // ambiguous_function
	ERRCODE_AMBIGUOUS_PARAMETER           = "42P08" // ambiguous_parameter
	ERRCODE_AMBIGUOUS_ALIAS               = "42P09" // ambiguous_alias
	ERRCODE_INVALID_COLUMN_REFERENCE      = "42P10" // invalid_column_reference
	ERRCODE_INVALID_COLUMN_DEFINITION     = "42611" // invalid_column_definition
	ERRCODE_INVALID_CURSOR_DEFINITION     = "42P11" // invalid_cursor_definition
	ERRCODE_INVALID_DATABASE_DEFINITION   = "42P12" // invalid_database_definition
	ERRCODE_INVALID_FUNCTION_DEFINITION   = "42P13" // invalid_function_definition
	ERRCODE_INVALID_PSTATEMENT_DEFINITION = "42P14" // invalid_prepared_statement_definition
	ERRCODE_INVALID_SCHEMA_DEFINITION     = "42P15" // invalid_schema_definition
	ERRCODE_INVALID_TABLE_DEFINITION      = "42P16" // invalid_table_definition
	ERRCODE_INVALID_OBJECT_DEFINITION     = "42P17" // invalid_object_definition

	// Section: Class 44 - WITH CHECK OPTION Violation.

	ERRCODE_WITH_CHECK_OPTION_VIOLATION = "44000" // with_check_option_violation

	// Section: Class 53 - Insufficient Resources.

	ERRCODE_INSUFFICIENT_RESOURCES       = "53000" // insufficient_resources
	ERRCODE_DISK_FULL                    = "53100" // disk_full
	ERRCODE_OUT_OF_MEMORY                = "53200" // out_of_memory
	ERRCODE_TOO_MANY_CONNECTIONS         = "53300" // too_many_connections
	ERRCODE_CONFIGURATION_LIMIT_EXCEEDED = "53400" // configuration_limit_exceeded

	// Section: Class 54 - Program Limit Exceeded.

	ERRCODE_PROGRAM_LIMIT_EXCEEDED = "54000" // program_limit_exceeded
	ERRCODE_STATEMENT_TOO_COMPLEX  = "54001" // statement_too_complex
	ERRCODE_TOO_MANY_COLUMNS       = "54011" // too_many_columns
	ERRCODE_TOO_MANY_ARGUMENTS     = "54023" // too_many_arguments

	// Section: Class 55 - Object Not In Prerequisite State.

	ERRCODE_OBJECT_NOT_IN_PREREQUISITE_STATE = "55000" // object_not_in_prerequisite_state
	ERRCODE_OBJECT_IN_USE                    = "55006" // object_in_use
	ERRCODE_CANT_CHANGE_RUNTIME_PARAM        = "55P02" // cant_change_runtime_param
	ERRCODE_LOCK_NOT_AVAILABLE               = "55P03" // lock_not_available
	ERRCODE_UNSAFE_NEW_ENUM_VALUE_USAGE      = "55P04" // unsafe_new_enum_value_usage

	// Section: Class 57 - Operator Intervention.

	ERRCODE_OPERATOR_INTERVENTION = "57000" // operator_intervention
	ERRCODE_QUERY_CANCELED        = "57014" // query_canceled
	ERRCODE_ADMIN_SHUTDOWN        = "57P01" // admin_shutdown
	ERRCODE_CRASH_SHUTDOWN        = "57P02" // crash_shutdown
	ERRCODE_CANNOT_CONNECT_NOW    = "57P03" // cannot_connect_now
	ERRCODE_DATABASE_DROPPED      = "57P04" // database_dropped
	ERRCODE_IDLE_SESSION_TIMEOUT  = "57P05" // idle_session_timeout

	// Section: Class 58 - System Error (errors external to PostgreSQL itself).

	ERRCODE_SYSTEM_ERROR   = "58000" // system_error
	ERRCODE_IO_ERROR       = "58030" // io_error
	ERRCODE_UNDEFINED_FILE = "58P01" // undefined_file
	ERRCODE_DUPLICATE_FILE = "58P02" // duplicate_file

	// Section: Class 72 - Snapshot Failure.
	ERRCODE_SNAPSHOT_TOO_OLD = "72000" // snapshot_too_old

	// Section: Class F0 - Configuration File Error.

	ERRCODE_CONFIG_FILE_ERROR = "F0000" // config_file_error
	ERRCODE_LOCK_FILE_EXISTS  = "F0001" // lock_file_exists

	// Section: Class HV - Foreign Data Wrapper Error (SQL/MED).

	ERRCODE_FDW_ERROR                                  = "HV000" // fdw_error
	ERRCODE_FDW_COLUMN_NAME_NOT_FOUND                  = "HV005" // fdw_column_name_not_found
	ERRCODE_FDW_DYNAMIC_PARAMETER_VALUE_NEEDED         = "HV002" // fdw_dynamic_parameter_value_needed
	ERRCODE_FDW_FUNCTION_SEQUENCE_ERROR                = "HV010" // fdw_function_sequence_error
	ERRCODE_FDW_INCONSISTENT_DESCRIPTOR_INFORMATION    = "HV021" // fdw_inconsistent_descriptor_information
	ERRCODE_FDW_INVALID_ATTRIBUTE_VALUE                = "HV024" // fdw_invalid_attribute_value
	ERRCODE_FDW_INVALID_COLUMN_NAME                    = "HV007" // fdw_invalid_column_name
	ERRCODE_FDW_INVALID_COLUMN_NUMBER                  = "HV008" // fdw_invalid_column_number
	ERRCODE_FDW_INVALID_DATA_TYPE                      = "HV004" // fdw_invalid_data_type
	ERRCODE_FDW_INVALID_DATA_TYPE_DESCRIPTORS          = "HV006" // fdw_invalid_data_type_descriptors
	ERRCODE_FDW_INVALID_DESCRIPTOR_FIELD_IDENTIFIER    = "HV091" // fdw_invalid_descriptor_field_identifier
	ERRCODE_FDW_INVALID_HANDLE                         = "HV00B" // fdw_invalid_handle
	ERRCODE_FDW_INVALID_OPTION_INDEX                   = "HV00C" // fdw_invalid_option_index
	ERRCODE_FDW_INVALID_OPTION_NAME                    = "HV00D" // fdw_invalid_option_name
	ERRCODE_FDW_INVALID_STRING_LENGTH_OR_BUFFER_LENGTH = "HV090" // fdw_invalid_string_length_or_buffer_length
	ERRCODE_FDW_INVALID_STRING_FORMAT                  = "HV00A" // fdw_invalid_string_format
	ERRCODE_FDW_INVALID_USE_OF_NULL_POINTER            = "HV009" // fdw_invalid_use_of_null_pointer
	ERRCODE_FDW_TOO_MANY_HANDLES                       = "HV014" // fdw_too_many_handles
	ERRCODE_FDW_OUT_OF_MEMORY                          = "HV001" // fdw_out_of_memory
	ERRCODE_FDW_NO_SCHEMAS                             = "HV00P" // fdw_no_schemas
	ERRCODE_FDW_OPTION_NAME_NOT_FOUND                  = "HV00J" // fdw_option_name_not_found
	ERRCODE_FDW_REPLY_HANDLE                           = "HV00K" // fdw_reply_handle
	ERRCODE_FDW_SCHEMA_NOT_FOUND                       = "HV00Q" // fdw_schema_not_found
	ERRCODE_FDW_TABLE_NOT_FOUND                        = "HV00R" // fdw_table_not_found
	ERRCODE_FDW_UNABLE_TO_CREATE_EXECUTION             = "HV00L" // fdw_unable_to_create_execution
	ERRCODE_FDW_UNABLE_TO_CREATE_REPLY                 = "HV00M" // fdw_unable_to_create_reply
	ERRCODE_FDW_UNABLE_TO_ESTABLISH_CONNECTION         = "HV00N" // fdw_unable_to_establish_connection

	// Section: Class P0 - PL/pgSQL Error.

	// (PostgreSQL-specific error class).
	ERRCODE_PLPGSQL_ERROR   = "P0000" // plpgsql_error
	ERRCODE_RAISE_EXCEPTION = "P0001" // raise_exception
	ERRCODE_NO_DATA_FOUND   = "P0002" // no_data_found
	ERRCODE_TOO_MANY_ROWS   = "P0003" // too_many_rows
	ERRCODE_ASSERT_FAILURE  = "P0004" // assert_failure

	// Section: Class XX - Internal Error.

	// this is for "can't-happen" conditions and software bugs (PostgreSQL-specific error class).
	ERRCODE_INTERNAL_ERROR  = "XX000" // internal_error
	ERRCODE_DATA_CORRUPTED  = "XX001" // data_corrupted
	ERRCODE_INDEX_CORRUPTED = "XX002" // index_corrupted
)
