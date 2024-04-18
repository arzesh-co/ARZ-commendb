package CommenDb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GetErrors(key string, account string, lang string, params map[string]string) *ResponseErrors {
	res := getError(key, account, lang, params)
	if res == nil {
		return nil
	}
	return res
}

func InvalidError() string {
	return `{
		"_id": "c016f23a-8924-4ebf-b80d-297fa1c2b9fc",
		"error_type": {
			"_id": "89df6266-2251-4323-96f7-9e0cda327562",
			"key": "input",
			"title": {
				"en": "input validation",
				"fa": "اعتبار سنجی ورودی"
			},
			"desc": {
				"en": "input value is not valid",
				"fa": "نامعتبر بودن اطلاعات ورودی"
			},
			"devTeamId": "",
			"url": "",
			"status": 1
		},
		"status_key": "REF.INVALIDATION_ERROR",
        "status_code": 422,
		"detail": {
			"en": "The submitted $information$ is not valid",
			"fa": "$information$ ارسالی معتبر نمی باشد"
		},
		"service_key": "",
		"title": {
			"en": "submitted",
			"fa": "ارسال اطلاعات"
		},
		"params": [
			{
				"key": "$information$",
				"default": {
					"en": "information",
					"fa": "اطلاعات"
				}
			}
		],
		"help_url": "",
		"meta_data": { },
		"status": 1
	}`
}

func ConnectionsError() string {
	return `{
		"_id": "647ca2d3-f97e-4331-8aba-046a97ec8682",
		"error_type": {
			"_id": "1ec55988-8412-45aa-be3d-1311a5c7acee",
			"key": "connection",
			"title": {
				"en": "connection",
				"fa": "ارتباطات"
			},
			"desc": {
				"en": "error to create connections",
				"fa": "خطا در ایجاد ارتباطات"
			},
			"devTeamId": "",
			"url": "",
			"status": 1
		},
		"status_key": "REF.CANNOT_CONNECT",
        "status_code": 503,
		"detail": {
			"fa": "خطا در ارتباط با $service$",
			"en": "error to connection with $service$"
		},
		"service_key": "",
		"title": {
			"en": "connections",
			"fa": "ارتباطات"
		},
		"params": [
			{
				"key": "$service$",
				"default": {
					"fa": "سرویس",
					"en": "service"
				}
			}
		],
		"help_url": "",
		"meta_data": { },
		"status": 1
	}`
}

func NotFoundError() string {
	return `{
		"_id": "6523900a-2783-48c7-8300-5a4a7e24058d",
		"error_type": {
			"_id": "c407cad6-a087-491c-af22-a5a2374354e3",
			"key": "find",
			"title": {
				"en": "find error",
				"fa": "خطا در یافتن"
			},
			"desc": {
				"en": "error to find entity",
				"fa": "خطا در یافتن موجودیت"
			},
			"devTeamId": "",
			"url": "",
			"status": 1
		},
		"status_key": "REF.NOT_FOUND",
        "status_code": 404,
		"detail": {
			"en": "$data$ can't find.",
			"fa": "$data$ یافت نشد"
		},
		"service_key": "",
		"title": {
			"en": "find data",
			"fa": "یافتن اطلاعات"
		},
		"params": [
			{
				"key": "$data$",
				"default": {
					"en": "data",
					"fa": "اطلاعات"
				}
			}
		],
		"help_url": "",
		"meta_data": null,
		"status": 1
	}`
}

func AlreadyInsertedError() string {
	return `{
		"_id": "7b4b7729-9b0b-4b75-b49e-a4f36384c0dc",
		"error_type": {
			"_id": "93053aac-de44-45ef-a838-00cd38d06f90",
			"key": "repeat",
			"title": {
				"fa": "تکرار",
				"en": "repeat"
			},
			"desc": {
				"en": "repeat to do Process",
				"fa": "تکرار در انجام یک فرایند"
			},
			"devTeamId": "",
			"url": "",
			"status": 1
		},
		"status_key": "REF.ALREADY_INSERTED",
        "status_code": 409,
		"detail": {
			"en": "The $information$ has been recorded in the past",
			"fa": "$information$ در گذشته ثبت شده است"
		},
		"service_key": "",
		"title": {
			"fa": "تکراری اطلاعات",
			"en": "repeat information"
		},
		"params": [
			{
				"key": "$information$",
				"default": {
					"en": "information",
					"fa": "اطلاعات"
				}
			}
		],
		"help_url": "",
		"meta_data": { },
		"status": 1
	}`
}

func AccessError() string {
	return `{
    "_id": "5a1a659f-eb23-4ef2-9762-80adc0c978ee",
    "error_type": {
        "_id": "9ee0f5c7-1848-424d-a269-ba5d4e653af5",
        "key": "access",
        "title": {
            "en": "access",
            "fa": "دسترسی"
        },
        "desc": {
            "en": "Insufficient access level",
            "fa": "کافی نبودن سطح دسترسی"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "REF.CANNOT_ACCESS",
    "status_code": 403,
    "detail": {
        "en": "Your access level is not enough",
        "fa": "سطح دسترسی شما کافی نمی باشد"
    },
    "service_key": "",
    "title": {
        "en": "access level",
        "fa": "سطح دسترسی"
    },
    "params": [ ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func InsertError() string {
	return `{
    "_id": "6682452f-f9e0-4065-bf2a-7e6a1880e638",
    "error_type": {
        "_id": "ed2729c7-b97f-40f0-9adc-554fd7ad7a67",
        "key": "write_in",
        "title": {
            "en": "write",
            "fa": "ثبت"
        },
        "desc": {
            "en": "error during write record",
            "fa": "خطا در ثبت"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "REF.CANNOT_INSERT",
    "status_code": 409
    "detail": {
        "fa": "خطا در ثبت $information$",
        "en": "Error recording $information$"
    },
    "service_key": "",
    "title": {
        "en": "recording information",
        "fa": "ثبت اطلاعات"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "en": "information",
                "fa": "اطلاعات"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
	}`
}

func ExpiredError() string {
	return `{
    "_id": "a82fee84-c423-4f40-b2b9-bca4939311f1",
    "error_type": {
        "_id": "ba338eac-f915-49b6-80c3-65568cdb1288",
        "key": "expiration",
        "title": {
            "en": "expiration",
            "fa": "انقضا"
        },
        "desc": {
            "fa": "منقضی شدن اطلاعات",
            "en": "expiration data value"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "REF.DATA_EXPIRED",
    "status_code": 419,
    "detail": {
        "en": "$information$ was expiraded",
        "fa": "$information$ منقضی شده است"
    },
    "service_key": "",
    "title": {
        "fa": "منقضی شده",
        "en": "expiraded"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "fa": "اطلاعات",
                "en": "information"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func ResponseError() string {
	return `{
    "_id": "2cfdbe2c-7925-435f-b34d-a67e2affbc59",
    "error_type": {
        "_id": "1ec55988-8412-45aa-be3d-1311a5c7acee",
        "key": "connection",
        "title": {
            "en": "connection",
            "fa": "ارتباطات"
        },
        "desc": {
            "en": "error to create connections",
            "fa": "خطا در ایجاد ارتباطات"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "REF.RESPONSE_ERROR",
    "status_code": 400,
    "detail": {
        "en": "$service$ Unresponsived",
        "fa": "عدم پاسخ گویی $service$"
    },
    "service_key": "",
    "title": {
        "en": "Unresponsived",
        "fa": "عدم پاسخ گویی"
    },
    "params": [
        {
            "key": "$service$",
            "default": {
                "en": "service",
                "fa": "سرویس"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func WrongInformationError() string {
	return `{
    "_id": "2e7e9848-0e4e-4054-ae52-465be60c6868",
    "error_type": {
        "_id": "89df6266-2251-4323-96f7-9e0cda327562",
        "key": "input",
        "title": {
            "en": "input validation",
            "fa": "اعتبار سنجی ورودی"
        },
        "desc": {
            "en": "input value is not valid",
            "fa": "نامعتبر بودن اطلاعات ورودی"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "REF.WRONG_INFORMATION",
    "status_code": 422,
    "detail": {
        "en": "$information$ is wrong.",
        "fa": "$information$ صحیح نمی باشد"
    },
    "service_key": "",
    "title": {
        "fa": "اطلاعات",
        "en": "information"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "en": "information",
                "fa": "اطلاعات"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func UpdateError() string {
	return `{
    "_id": "751ad308-8378-43b8-9590-99142d920670",
    "error_type": {
        "_id": "ed2729c7-b97f-40f0-9adc-554fd7ad7a67",
        "key": "write_in",
        "title": {
            "en": "write",
            "fa": "ثبت"
        },
        "desc": {
            "en": "error during write record",
            "fa": "خطا در ثبت"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "REF.CANNOT_UPDATE",
    "status_code": 400,
    "detail": {
        "en": "error to update $information$",
        "fa": "خطا در بروزرسانی $information$"
    },
    "service_key": "",
    "title": {
        "en": "information",
        "fa": "اطلاعات"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "en": "information",
                "fa": "اطلاعات"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func RequiredError() string {
	return `{
    "_id": "cfc36d5b-9d63-4e57-add1-fe0b16732b4d",
    "error_type": {
        "_id": "89df6266-2251-4323-96f7-9e0cda327562",
        "key": "input",
        "title": {
            "en": "input validation",
            "fa": "اعتبار سنجی ورودی"
        },
        "desc": {
            "en": "input value is not valid",
            "fa": "نامعتبر بودن اطلاعات ورودی"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "required",
    "status_code": 400,
    "detail": {
        "en": "$information$ is required",
        "fa": "$information$ الزامی می باشد"
    },
    "service_key": "",
    "title": {
        "en": "required",
        "fa": "الزامی"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "en": "information",
                "fa": "اطلاعات"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func AlphaNumError() string {
	return `{
    "_id": "655194cc-88c7-4dea-a7e7-000acade744f",
    "error_type": {
        "_id": "89df6266-2251-4323-96f7-9e0cda327562",
        "key": "input",
        "title": {
            "en": "input validation",
            "fa": "اعتبار سنجی ورودی"
        },
        "desc": {
            "en": "input value is not valid",
            "fa": "نامعتبر بودن اطلاعات ورودی"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "alphaNum",
    "status_code": 400,
    "detail": {
        "en": "$information$ must be only alphanumerics.",
        "fa": "$information$ باید حروف و یا اعداد باشد"
    },
    "service_key": "",
    "title": {
        "en": "alphaNum",
        "fa": "اعداد و حروف"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "en": "information",
                "fa": "اطلاعات"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func MinLengthError() string {
	return `{
    "_id": "6072d524-f8b2-42e5-9b4b-02d2da1c14a4",
    "error_type": {
        "_id": "89df6266-2251-4323-96f7-9e0cda327562",
        "key": "input",
        "title": {
            "en": "input validation",
            "fa": "اعتبار سنجی ورودی"
        },
        "desc": {
            "en": "input value is not valid",
            "fa": "نامعتبر بودن اطلاعات ورودی"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "minLength",
    "status_code": 400,
    "detail": {
        "en": "$information$ length must be more.",
        "fa": "تعداد کاراکتر $information$ کمتر از حد مجاز است"
    },
    "service_key": "",
    "title": {
        "en": "minLength",
        "fa": "حداقل کاراکتر"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "en": "information",
                "fa": "اطلاعات"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func MaxLengthError() string {
	return `{
    "_id": "77675b9f-c592-4ca0-9607-fa91ef40f2f0",
    "error_type": {
        "_id": "89df6266-2251-4323-96f7-9e0cda327562",
        "key": "input",
        "title": {
            "en": "input validation",
            "fa": "اعتبار سنجی ورودی"
        },
        "desc": {
            "en": "input value is not valid",
            "fa": "نامعتبر بودن اطلاعات ورودی"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "maxLength",
    "detail": {
        "en": "$information$ length is more than maxLength",
        "fa": "تعداد کاراکتر $information$ بیشتر از حد مجاز است"
    },
    "service_key": "",
    "title": {
        "en": "maxLength",
        "fa": "حداکثر کاراکتر"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "en": "information",
                "fa": "اطلاعات"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func MinValueError() string {
	return `{
    "_id": "c6b7ce4a-9e5f-429e-8626-6126902204f2",
    "error_type": {
        "_id": "89df6266-2251-4323-96f7-9e0cda327562",
        "key": "input",
        "title": {
            "en": "input validation",
            "fa": "اعتبار سنجی ورودی"
        },
        "desc": {
            "en": "input value is not valid",
            "fa": "نامعتبر بودن اطلاعات ورودی"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "minValue",
    "status_code": 400,
    "detail": {
        "en": "$information$ length is less than minValue",
        "fa": "مقدار $information$ کمتر از حد مجاز است"
    },
    "service_key": "",
    "title": {
        "en": "minValue",
        "fa": "حداقل مقدار"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "en": "information",
                "fa": "اطلاعات"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func MaxValueError() string {
	return `{
    "_id": "d3a3100b-fed7-4ad3-a314-8d68c57790ba",
    "error_type": {
        "_id": "89df6266-2251-4323-96f7-9e0cda327562",
        "key": "input",
        "title": {
            "en": "input validation",
            "fa": "اعتبار سنجی ورودی"
        },
        "desc": {
            "en": "input value is not valid",
            "fa": "نامعتبر بودن اطلاعات ورودی"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "maxValue",
    "status_code": 400,
    "detail": {
        "en": "$information$ length is more than maxValue",
        "fa": "مقدار $information$ بیشتر از حد مجاز است"
    },
    "service_key": "",
    "title": {
        "en": "maxValue",
        "fa": "حداکثر مقدار"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "en": "information",
                "fa": "اطلاعات"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func PhoneError() string {
	return `{
    "_id": "8750ce6d-87c6-42ab-8cc0-aac898b3dfdb",
    "error_type": {
        "_id": "89df6266-2251-4323-96f7-9e0cda327562",
        "key": "input",
        "title": {
            "en": "input validation",
            "fa": "اعتبار سنجی ورودی"
        },
        "desc": {
            "en": "input value is not valid",
            "fa": "نامعتبر بودن اطلاعات ورودی"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "phone",
    "status_code": 400,
    "detail": {
        "en": "$information$ is not phone format.",
        "fa": "$information$ با فرمت تلفن همخوانی ندارد"
    },
    "service_key": "",
    "title": {
        "en": "phone",
        "fa": "تلفن"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "en": "information",
                "fa": "اطلاعات"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func EmailError() string {
	return `{
    "_id": "ece92a69-f250-4d46-a7ec-ae91e4b0c808",
    "error_type": {
        "_id": "89df6266-2251-4323-96f7-9e0cda327562",
        "key": "input",
        "title": {
            "en": "input validation",
            "fa": "اعتبار سنجی ورودی"
        },
        "desc": {
            "en": "input value is not valid",
            "fa": "نامعتبر بودن اطلاعات ورودی"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "email",
    "status_code": 400,
    "detail": {
        "en": "$information$ is not email format.",
        "fa": "$information$ با فرمت ایمیل همخوانی ندارد"
    },
    "service_key": "",
    "title": {
        "en": "email",
        "fa": "ایمیل"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "en": "information",
                "fa": "اطلاعات"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func NumericError() string {
	return `{
    "_id": "b2f87e73-fb54-4674-bc66-a627f619579e",
    "error_type": {
        "_id": "89df6266-2251-4323-96f7-9e0cda327562",
        "key": "input",
        "title": {
            "en": "input validation",
            "fa": "اعتبار سنجی ورودی"
        },
        "desc": {
            "en": "input value is not valid",
            "fa": "نامعتبر بودن اطلاعات ورودی"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "numeric",
    "status_code": 400,
    "detail": {
        "en": "$information$ is not numeric format.",
        "fa": "$information$ با فرمت عددی همخوانی ندارد"
    },
    "service_key": "",
    "title": {
        "en": "numeric",
        "fa": "اعداد"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "fa": "اطلاعات",
                "en": "information"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func IntegerError() string {
	return `{
    "_id": "e4fd80b6-6d5c-46d0-a01e-268a51867c87",
    "error_type": {
        "_id": "89df6266-2251-4323-96f7-9e0cda327562",
        "key": "input",
        "title": {
            "en": "input validation",
            "fa": "اعتبار سنجی ورودی"
        },
        "desc": {
            "en": "input value is not valid",
            "fa": "نامعتبر بودن اطلاعات ورودی"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "number",
    "status_code": 400,
    "detail": {
        "en": "$information$ is not number",
        "fa": "$information$ عدد صحیح نمی باشد"
    },
    "service_key": "",
    "title": {
        "en": "integer",
        "fa": "عدد صحیح"
    },
    "params": [
        {
            "key": "$information$",
            "default": {
                "fa": "اطلاعات",
                "en": "information"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func DependencyError() string {
	return `{
    "_id": "30cc6d87-cf31-403a-bebd-f816d4dc929c",
    "error_type": {
        "_id": "3a9698a3-e290-4a6a-b325-ac501ff8e3a1",
        "key": "dependence",
        "title": {
            "fa": "وابستگی ها",
            "en": "dependence"
        },
        "desc": {
            "en": "dependence error",
            "fa": "خطا در وابستگی های موجودیت"
        },
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "REF.DATA_IS_DEPENDENT",
    "status_code": 424,
    "detail": {
        "en": "$data$ is dependent",
        "fa": "$data$ دارای وابستگی می باشد"
    },
    "service_key": "",
    "title": {
        "fa": "وابستگی",
        "en": "dependency"
    },
    "params": [
        {
            "key": "$data$",
            "default": {
                "en": "data",
                "fa": "داده"
            }
        }
    ],
    "help_url": "",
    "meta_data": { },
    "status": 1
}`
}

func RemoveError() string {
	return `{
    "_id": "18a41710-eedc-4d72-8c55-d6d607c26b1e",
    "error_type": {
        "_id": "ade08ab8-0ebf-4bcd-ba03-0f3a90d0476f",
        "key": "delete",
        "title": {
            "en": "delete",
            "fa": "خطا در حذف"
        },
        "desc": {
            "en": "error when delete something",
            "fa": "خطا هنگام حذف داده"
        },
        "service_key": "",
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "REF.CANNOT_REMOVE",
    "status_code": 400,
    "detail": {
        "en": "$data$ can't delete",
        "fa": "مشکلی در حذف بوجود آمده"
    },
    "service_key": "",
    "title": {
        "en": "delete",
        "fa": "خطا در حذف"
    },
    "params": [
        {
            "key": "$data$",
            "default": {
                "en": "data",
                "fa": "اطلاعات"
            }
        }
    ],
    "help_url": "",
    "meta_data": null,
    "status": 1
}`
}

func UnknownError() string {
	return `{
    "_id": "d998a8ea-99ea-43f9-b264-ba74da706dd1",
    "error_type": {
        "_id": "",
        "key": "service",
        "title": null,
        "desc": null,
        "service_key": "",
        "devTeamId": "",
        "url": "",
        "status": 0
    },
    "status_key": "REF.SERVICE_UNKNOWN_ERROR",
    "status_code": 400,
    "detail": {
        "en": "service has unknown error.",
        "fa": "سرویس با مشکل مواجه شده است لطفا مجددا تلاش فرمایید"
    },
    "service_key": "",
    "title": {
        "en": "service error",
        "fa": "خطای سرویس"
    },
    "params": [],
    "help_url": "",
    "meta_data": null,
    "status": 0
}`
}

func SettingError() string {
	return `{
    "_id": "421e5fcd-0c85-4963-8c61-489b457f9378",
    "error_type": {
        "_id": "dc0da5e5-7bbc-4456-8d1b-315bc51f4568",
        "key": "service",
        "title": {
            "en": "service error",
            "fa": "خطای سرویس"
        },
        "desc": {
            "en": "service error",
            "fa": "خطا هنگام انجام سرویس"
        },
        "service_key": "",
        "devTeamId": "",
        "url": "",
        "status": 1
    },
    "status_key": "REF.CONFIG_SETTING_NOT_FOUND",
    "status_code": 400,
    "detail": {
        "en": "config account not found",
        "fa": "تنظیمات اکانت یافت نشد"
    },
    "service_key": "",
    "title": {
        "en": "config not found",
        "fa": "موجود نبودن تنظیمات"
    },
    "params": [],
    "help_url": "",
    "meta_data": null,
    "status": 1
}`
}

var ErrorKeyMap = map[string]string{
	"REF.INVALIDATION_ERROR":       InvalidError(),
	"REF.CANNOT_CONNECT":           ConnectionsError(),
	"REF.NOT_FOUND":                NotFoundError(),
	"REF.ALREADY_INSERTED":         AlreadyInsertedError(),
	"REF.CANNOT_ACCESS":            AccessError(),
	"REF.CANNOT_INSERT":            InsertError(),
	"REF.DATA_EXPIRED":             ExpiredError(),
	"REF.RESPONSE_ERROR":           ResponseError(),
	"REF.WRONG_INFORMATION":        WrongInformationError(),
	"REF.CANNOT_UPDATE":            UpdateError(),
	"required":                     RequiredError(),
	"alphaNum":                     AlphaNumError(),
	"minLength":                    MinLengthError(),
	"maxLength":                    MaxLengthError(),
	"minValue":                     MinValueError(),
	"maxValue":                     MaxValueError(),
	"phone":                        PhoneError(),
	"email":                        EmailError(),
	"numeric":                      NumericError(),
	"number":                       IntegerError(),
	"REF.DATA_IS_DEPENDENT":        DependencyError(),
	"REF.CANNOT_REMOVE":            RemoveError(),
	"REF.SERVICE_UNKNOWN_ERROR":    UnknownError(),
	"REF.CONFIG_SETTING_NOT_FOUND": SettingError(),
}

type ErrorType struct {
	Id         string            `json:"_id" bson:"_id"`
	Key        string            `json:"key" bson:"key"`
	Title      map[string]string `json:"title" bson:"title"`
	Desc       map[string]string `json:"desc" bson:"desc"`
	ServiceKey string            `json:"service_key" bson:"service_key"`
	DevTeamId  string            `json:"dev_team_uuid" bson:"devTeamId"`
	Url        string            `json:"url" bson:"url"`
	Status     int8              `json:"status" bson:"status"`
}

type ErrorParam struct {
	Key          string            `json:"key" bson:"key"`
	DefaultValue map[string]string `json:"default" bson:"default"`
}

type Errors struct {
	Id         string            `json:"_id" bson:"_id"`
	ErrorType  ErrorType         `json:"error_type" bson:"error_type"`
	StatusKey  string            `json:"status_key" bson:"status_key"`
	StatusCode int               `json:"status_code" bson:"status_code"`
	Detail     map[string]string `json:"detail" bson:"detail"`
	ServiceKey string            `json:"service_key" bson:"service_key"`
	Title      map[string]string `json:"title" bson:"title"`
	Params     []ErrorParam      `json:"params" bson:"params"`
	HelpUrl    string            `json:"help_url" bson:"help_url"`
	MetaData   map[string]any    `json:"meta_data" bson:"meta_data"`
	Status     int8              `json:"status" bson:"status"`
}

type Entities struct {
	Id         string            `json:"_id" bson:"_id"`
	EntityName string            `json:"entity_name" bson:"entity_name"`
	Title      map[string]string `json:"title" bson:"title"`
}

func FindError(key string) *Errors {
	strErr, ok := ErrorKeyMap[key]
	if !ok {
		return GetCustomErrorService(key)
	}

	Err := &Errors{}
	err := json.Unmarshal([]byte(strErr), Err)
	if err != nil {
		return nil
	}
	return Err
}

func convertErrorToResponseErr(Err *Errors, lang string) *ResponseErrors {
	res := &ResponseErrors{}
	res.ErrorType.Url = Err.ErrorType.Url
	res.ErrorType.Title = Err.ErrorType.Title[lang]
	res.ErrorType.Desc = Err.ErrorType.Desc[lang]
	res.Detail = Err.Detail[lang]
	res.HelpUrl = Err.HelpUrl
	res.Title = Err.Title[lang]
	res.StatusKey = Err.StatusKey
	return res
}

func setParamsToResponseErr(res *ResponseErrors, DefaultParams []ErrorParam, params map[string]string, lang string) *ResponseErrors {
	for _, p := range DefaultParams {

		paramVal, ok := params[p.Key]
		if !ok {
			res.Detail = strings.Replace(res.Detail, p.Key, p.DefaultValue[p.Key], -1)
			continue
		}

		res.Detail = strings.Replace(res.Detail, p.Key, paramVal, -1)
	}
	return res
}

type Param struct {
	Key     string            `json:"key" bson:"key"`
	Default map[string]string `json:"default" bson:"default"`
}

type ResErrorService struct {
	Data *Errors `json:"data"`
}

func GetCustomErrorService(errorKey string) *Errors {
	req, err := http.NewRequest("GET", os.Getenv("error_handling")+"/api/error-handling/errors/key/"+errorKey, nil)
	if err != nil {
		return nil
	}

	client := &http.Client{
		Transport: &http.Transport{},
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error is :", err.Error())
		return nil
	}
	req.Close = true
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error is :", err.Error())
		return nil
	}

	Info := &ResErrorService{}
	err = json.Unmarshal(body, Info)
	if err != nil {
		fmt.Println("error is :", err.Error())
		return nil
	}

	return Info.Data
}
