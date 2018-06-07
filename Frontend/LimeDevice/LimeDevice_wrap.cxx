/* ----------------------------------------------------------------------------
 * This file was automatically generated by SWIG (http://www.swig.org).
 * Version 3.0.12
 *
 * This file is not intended to be easily readable and contains a number of
 * coding conventions designed to improve portability and efficiency. Do not make
 * changes to this file unless you know what you are doing--modify the SWIG
 * interface file instead.
 * ----------------------------------------------------------------------------- */

// source: Frontend/LimeDevice/LimeDevice.i

#define SWIGMODULE LimeDevice
#define SWIG_DIRECTORS

#ifdef __cplusplus
/* SwigValueWrapper is described in swig.swg */
template<typename T> class SwigValueWrapper {
  struct SwigMovePointer {
    T *ptr;
    SwigMovePointer(T *p) : ptr(p) { }
    ~SwigMovePointer() { delete ptr; }
    SwigMovePointer& operator=(SwigMovePointer& rhs) { T* oldptr = ptr; ptr = 0; delete oldptr; ptr = rhs.ptr; rhs.ptr = 0; return *this; }
  } pointer;
  SwigValueWrapper& operator=(const SwigValueWrapper<T>& rhs);
  SwigValueWrapper(const SwigValueWrapper<T>& rhs);
public:
  SwigValueWrapper() : pointer(0) { }
  SwigValueWrapper& operator=(const T& t) { SwigMovePointer tmp(new T(t)); pointer = tmp; return *this; }
  operator T&() const { return *pointer.ptr; }
  T *operator&() { return pointer.ptr; }
};

template <typename T> T SwigValueInit() {
  return T();
}
#endif

/* -----------------------------------------------------------------------------
 *  This section contains generic SWIG labels for method/variable
 *  declarations/attributes, and other compiler dependent labels.
 * ----------------------------------------------------------------------------- */

/* template workaround for compilers that cannot correctly implement the C++ standard */
#ifndef SWIGTEMPLATEDISAMBIGUATOR
# if defined(__SUNPRO_CC) && (__SUNPRO_CC <= 0x560)
#  define SWIGTEMPLATEDISAMBIGUATOR template
# elif defined(__HP_aCC)
/* Needed even with `aCC -AA' when `aCC -V' reports HP ANSI C++ B3910B A.03.55 */
/* If we find a maximum version that requires this, the test would be __HP_aCC <= 35500 for A.03.55 */
#  define SWIGTEMPLATEDISAMBIGUATOR template
# else
#  define SWIGTEMPLATEDISAMBIGUATOR
# endif
#endif

/* inline attribute */
#ifndef SWIGINLINE
# if defined(__cplusplus) || (defined(__GNUC__) && !defined(__STRICT_ANSI__))
#   define SWIGINLINE inline
# else
#   define SWIGINLINE
# endif
#endif

/* attribute recognised by some compilers to avoid 'unused' warnings */
#ifndef SWIGUNUSED
# if defined(__GNUC__)
#   if !(defined(__cplusplus)) || (__GNUC__ > 3 || (__GNUC__ == 3 && __GNUC_MINOR__ >= 4))
#     define SWIGUNUSED __attribute__ ((__unused__))
#   else
#     define SWIGUNUSED
#   endif
# elif defined(__ICC)
#   define SWIGUNUSED __attribute__ ((__unused__))
# else
#   define SWIGUNUSED
# endif
#endif

#ifndef SWIG_MSC_UNSUPPRESS_4505
# if defined(_MSC_VER)
#   pragma warning(disable : 4505) /* unreferenced local function has been removed */
# endif
#endif

#ifndef SWIGUNUSEDPARM
# ifdef __cplusplus
#   define SWIGUNUSEDPARM(p)
# else
#   define SWIGUNUSEDPARM(p) p SWIGUNUSED
# endif
#endif

/* internal SWIG method */
#ifndef SWIGINTERN
# define SWIGINTERN static SWIGUNUSED
#endif

/* internal inline SWIG method */
#ifndef SWIGINTERNINLINE
# define SWIGINTERNINLINE SWIGINTERN SWIGINLINE
#endif

/* exporting methods */
#if defined(__GNUC__)
#  if (__GNUC__ >= 4) || (__GNUC__ == 3 && __GNUC_MINOR__ >= 4)
#    ifndef GCC_HASCLASSVISIBILITY
#      define GCC_HASCLASSVISIBILITY
#    endif
#  endif
#endif

#ifndef SWIGEXPORT
# if defined(_WIN32) || defined(__WIN32__) || defined(__CYGWIN__)
#   if defined(STATIC_LINKED)
#     define SWIGEXPORT
#   else
#     define SWIGEXPORT __declspec(dllexport)
#   endif
# else
#   if defined(__GNUC__) && defined(GCC_HASCLASSVISIBILITY)
#     define SWIGEXPORT __attribute__ ((visibility("default")))
#   else
#     define SWIGEXPORT
#   endif
# endif
#endif

/* calling conventions for Windows */
#ifndef SWIGSTDCALL
# if defined(_WIN32) || defined(__WIN32__) || defined(__CYGWIN__)
#   define SWIGSTDCALL __stdcall
# else
#   define SWIGSTDCALL
# endif
#endif

/* Deal with Microsoft's attempt at deprecating C standard runtime functions */
#if !defined(SWIG_NO_CRT_SECURE_NO_DEPRECATE) && defined(_MSC_VER) && !defined(_CRT_SECURE_NO_DEPRECATE)
# define _CRT_SECURE_NO_DEPRECATE
#endif

/* Deal with Microsoft's attempt at deprecating methods in the standard C++ library */
#if !defined(SWIG_NO_SCL_SECURE_NO_DEPRECATE) && defined(_MSC_VER) && !defined(_SCL_SECURE_NO_DEPRECATE)
# define _SCL_SECURE_NO_DEPRECATE
#endif

/* Deal with Apple's deprecated 'AssertMacros.h' from Carbon-framework */
#if defined(__APPLE__) && !defined(__ASSERT_MACROS_DEFINE_VERSIONS_WITHOUT_UNDERSCORES)
# define __ASSERT_MACROS_DEFINE_VERSIONS_WITHOUT_UNDERSCORES 0
#endif

/* Intel's compiler complains if a variable which was never initialised is
 * cast to void, which is a common idiom which we use to indicate that we
 * are aware a variable isn't used.  So we just silence that warning.
 * See: https://github.com/swig/swig/issues/192 for more discussion.
 */
#ifdef __INTEL_COMPILER
# pragma warning disable 592
#endif


#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>



typedef long long intgo;
typedef unsigned long long uintgo;


# if !defined(__clang__) && (defined(__i386__) || defined(__x86_64__))
#   define SWIGSTRUCTPACKED __attribute__((__packed__, __gcc_struct__))
# else
#   define SWIGSTRUCTPACKED __attribute__((__packed__))
# endif



typedef struct { char *p; intgo n; } _gostring_;
typedef struct { void* array; intgo len; intgo cap; } _goslice_;




#define swiggo_size_assert_eq(x, y, name) typedef char name[(x-y)*(x-y)*-2+1];
#define swiggo_size_assert(t, n) swiggo_size_assert_eq(sizeof(t), n, swiggo_sizeof_##t##_is_not_##n)

swiggo_size_assert(char, 1)
swiggo_size_assert(short, 2)
swiggo_size_assert(int, 4)
typedef long long swiggo_long_long;
swiggo_size_assert(swiggo_long_long, 8)
swiggo_size_assert(float, 4)
swiggo_size_assert(double, 8)

#ifdef __cplusplus
extern "C" {
#endif
extern void crosscall2(void (*fn)(void *, int), void *, int);
extern char* _cgo_topofstack(void) __attribute__ ((weak));
extern void _cgo_allocate(void *, int);
extern void _cgo_panic(void *, int);
#ifdef __cplusplus
}
#endif

static char *_swig_topofstack() {
  if (_cgo_topofstack) {
    return _cgo_topofstack();
  } else {
    return 0;
  }
}

static void _swig_gopanic(const char *p) {
  struct {
    const char *p;
  } SWIGSTRUCTPACKED a;
  a.p = p;
  crosscall2(_cgo_panic, &a, (int) sizeof a);
}




#define SWIG_contract_assert(expr, msg) \
  if (!(expr)) { _swig_gopanic(msg); } else


#define SWIG_exception(code, msg) _swig_gopanic(msg)


static _gostring_ Swig_AllocateString(const char *p, size_t l) {
  _gostring_ ret;
  ret.p = (char*)malloc(l);
  memcpy(ret.p, p, l);
  ret.n = l;
  return ret;
}


static void Swig_free(void* p) {
  free(p);
}

static void* Swig_malloc(int c) {
  return malloc(c);
}


#include "LimeDevice.h"


#include <stdint.h>		// Use the C99 official header


#include <typeinfo>
#include <stdexcept>


#include <string>


#include <vector>
#include <stdexcept>


#include <map>
#include <algorithm>
#include <stdexcept>


#include <utility>

SWIGINTERN std::vector< unsigned int >::const_reference std_vector_Sl_uint32_t_Sg__get(std::vector< uint32_t > *self,int i){
                int size = int(self->size());
                if (i>=0 && i<size)
                    return (*self)[i];
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN void std_vector_Sl_uint32_t_Sg__set(std::vector< uint32_t > *self,int i,std::vector< unsigned int >::value_type const &val){
                int size = int(self->size());
                if (i>=0 && i<size)
                    (*self)[i] = val;
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN std::vector< float >::const_reference std_vector_Sl_float_Sg__get(std::vector< float > *self,int i){
                int size = int(self->size());
                if (i>=0 && i<size)
                    return (*self)[i];
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN void std_vector_Sl_float_Sg__set(std::vector< float > *self,int i,std::vector< float >::value_type const &val){
                int size = int(self->size());
                if (i>=0 && i<size)
                    (*self)[i] = val;
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN std::vector< short >::const_reference std_vector_Sl_int16_t_Sg__get(std::vector< int16_t > *self,int i){
                int size = int(self->size());
                if (i>=0 && i<size)
                    return (*self)[i];
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN void std_vector_Sl_int16_t_Sg__set(std::vector< int16_t > *self,int i,std::vector< short >::value_type const &val){
                int size = int(self->size());
                if (i>=0 && i<size)
                    (*self)[i] = val;
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN std::vector< signed char >::const_reference std_vector_Sl_int8_t_Sg__get(std::vector< int8_t > *self,int i){
                int size = int(self->size());
                if (i>=0 && i<size)
                    return (*self)[i];
                else
                    throw std::out_of_range("vector index out of range");
            }
SWIGINTERN void std_vector_Sl_int8_t_Sg__set(std::vector< int8_t > *self,int i,std::vector< signed char >::value_type const &val){
                int size = int(self->size());
                if (i>=0 && i<size)
                    (*self)[i] = val;
                else
                    throw std::out_of_range("vector index out of range");
            }

// C++ director class methods.
#include "LimeDevice_wrap.h"


#include <map>

namespace {
  struct GCItem {
    virtual ~GCItem() {}
  };

  struct GCItem_var {
    GCItem_var(GCItem *item = 0) : _item(item) {
    }

    GCItem_var& operator=(GCItem *item) {
      GCItem *tmp = _item;
      _item = item;
      delete tmp;
      return *this;
    }

    ~GCItem_var() {
      delete _item;
    }

    GCItem* operator->() {
      return _item;
    }

    private:
      GCItem *_item;
  };

  template <typename Type>
  struct GCItem_T : GCItem {
    GCItem_T(Type *ptr) : _ptr(ptr) {
    }

    virtual ~GCItem_T() {
      delete _ptr;
    }

  private:
    Type *_ptr;
  };
}

class Swig_memory {
public:
  template <typename Type>
  void swig_acquire_pointer(Type* vptr) {
    if (vptr) {
      swig_owner[vptr] = new GCItem_T<Type>(vptr);
    }
  }
private:
  typedef std::map<void *, GCItem_var> swig_ownership_map;
  swig_ownership_map swig_owner;
};

template <typename Type>
static void swig_acquire_pointer(Swig_memory** pmem, Type* ptr) {
  if (!pmem) {
    *pmem = new Swig_memory;
  }
  (*pmem)->swig_acquire_pointer(ptr);
}

SwigDirector_LimeCallback::SwigDirector_LimeCallback(int swig_p)
    : GoDeviceCallback(),
      go_val(swig_p), swig_mem(0)
{ }

extern "C" void Swig_DirectorLimeCallback_callback_cbFloatIQ_LimeDevice_f36a06c4aed6a5ac(int, void *arg2, intgo arg3);
void SwigDirector_LimeCallback::cbFloatIQ(void *data, int length) {
  void *swig_arg2;
  intgo swig_arg3;
  
  *(void **)&swig_arg2 = (void *)data; 
  swig_arg3 = (int)length; 
  Swig_DirectorLimeCallback_callback_cbFloatIQ_LimeDevice_f36a06c4aed6a5ac(go_val, swig_arg2, swig_arg3);
}

extern "C" void Swig_DirectorLimeCallback_callback_cbS16IQ_LimeDevice_f36a06c4aed6a5ac(int, short *arg2, intgo arg3);
void SwigDirector_LimeCallback::cbS16IQ(int16_t *data, int length) {
  short *swig_arg2;
  intgo swig_arg3;
  
  *(int16_t **)&swig_arg2 = (int16_t *)data; 
  swig_arg3 = (int)length; 
  Swig_DirectorLimeCallback_callback_cbS16IQ_LimeDevice_f36a06c4aed6a5ac(go_val, swig_arg2, swig_arg3);
}

extern "C" void Swig_DirectorLimeCallback_callback_cbS8IQ_LimeDevice_f36a06c4aed6a5ac(int, char *arg2, intgo arg3);
void SwigDirector_LimeCallback::cbS8IQ(int8_t *data, int length) {
  char *swig_arg2;
  intgo swig_arg3;
  
  *(int8_t **)&swig_arg2 = (int8_t *)data; 
  swig_arg3 = (int)length; 
  Swig_DirectorLimeCallback_callback_cbS8IQ_LimeDevice_f36a06c4aed6a5ac(go_val, swig_arg2, swig_arg3);
}

extern "C" void Swiggo_DeleteDirector_LimeCallback_LimeDevice_f36a06c4aed6a5ac(intgo);
SwigDirector_LimeCallback::~SwigDirector_LimeCallback()
{
  Swiggo_DeleteDirector_LimeCallback_LimeDevice_f36a06c4aed6a5ac(go_val);
  delete swig_mem;
}

#ifdef __cplusplus
extern "C" {
#endif

void _wrap_Swig_free_LimeDevice_f36a06c4aed6a5ac(void *_swig_go_0) {
  void *arg1 = (void *) 0 ;
  
  arg1 = *(void **)&_swig_go_0; 
  
  Swig_free(arg1);
  
}


void *_wrap_Swig_malloc_LimeDevice_f36a06c4aed6a5ac(intgo _swig_go_0) {
  int arg1 ;
  void *result = 0 ;
  void *_swig_go_result;
  
  arg1 = (int)_swig_go_0; 
  
  result = (void *)Swig_malloc(arg1);
  *(void **)&_swig_go_result = (void *)result; 
  return _swig_go_result;
}


GoDeviceCallback *_wrap__swig_NewDirectorLimeCallbackLimeCallback_LimeDevice_f36a06c4aed6a5ac(intgo _swig_go_0) {
  int arg1 ;
  GoDeviceCallback *result = 0 ;
  GoDeviceCallback *_swig_go_result;
  
  arg1 = (int)_swig_go_0; 
  
  result = new SwigDirector_LimeCallback(arg1);
  *(GoDeviceCallback **)&_swig_go_result = (GoDeviceCallback *)result; 
  return _swig_go_result;
}


void _wrap__swig_DirectorLimeCallback_upcall_CbFloatIQ_LimeDevice_f36a06c4aed6a5ac(SwigDirector_LimeCallback *_swig_go_0, void *_swig_go_1, intgo _swig_go_2) {
  SwigDirector_LimeCallback *arg1 = (SwigDirector_LimeCallback *) 0 ;
  void *arg2 = (void *) 0 ;
  int arg3 ;
  
  arg1 = *(SwigDirector_LimeCallback **)&_swig_go_0; 
  arg2 = *(void **)&_swig_go_1; 
  arg3 = (int)_swig_go_2; 
  
  arg1->_swig_upcall_cbFloatIQ(arg2, arg3);
  
}


void _wrap__swig_DirectorLimeCallback_upcall_CbS16IQ_LimeDevice_f36a06c4aed6a5ac(SwigDirector_LimeCallback *_swig_go_0, short *_swig_go_1, intgo _swig_go_2) {
  SwigDirector_LimeCallback *arg1 = (SwigDirector_LimeCallback *) 0 ;
  int16_t *arg2 = (int16_t *) 0 ;
  int arg3 ;
  
  arg1 = *(SwigDirector_LimeCallback **)&_swig_go_0; 
  arg2 = *(int16_t **)&_swig_go_1; 
  arg3 = (int)_swig_go_2; 
  
  arg1->_swig_upcall_cbS16IQ(arg2, arg3);
  
}


void _wrap__swig_DirectorLimeCallback_upcall_CbS8IQ_LimeDevice_f36a06c4aed6a5ac(SwigDirector_LimeCallback *_swig_go_0, char *_swig_go_1, intgo _swig_go_2) {
  SwigDirector_LimeCallback *arg1 = (SwigDirector_LimeCallback *) 0 ;
  int8_t *arg2 = (int8_t *) 0 ;
  int arg3 ;
  
  arg1 = *(SwigDirector_LimeCallback **)&_swig_go_0; 
  arg2 = *(int8_t **)&_swig_go_1; 
  arg3 = (int)_swig_go_2; 
  
  arg1->_swig_upcall_cbS8IQ(arg2, arg3);
  
}


void _wrap_DeleteDirectorLimeCallback_LimeDevice_f36a06c4aed6a5ac(GoDeviceCallback *_swig_go_0) {
  GoDeviceCallback *arg1 = (GoDeviceCallback *) 0 ;
  
  arg1 = *(GoDeviceCallback **)&_swig_go_0; 
  
  delete arg1;
  
}


void _wrap_LimeCallback_cbFloatIQ_LimeDevice_f36a06c4aed6a5ac(GoDeviceCallback *_swig_go_0, void *_swig_go_1, intgo _swig_go_2) {
  GoDeviceCallback *arg1 = (GoDeviceCallback *) 0 ;
  void *arg2 = (void *) 0 ;
  int arg3 ;
  
  arg1 = *(GoDeviceCallback **)&_swig_go_0; 
  arg2 = *(void **)&_swig_go_1; 
  arg3 = (int)_swig_go_2; 
  
  (arg1)->cbFloatIQ(arg2,arg3);
  
}


void _wrap_LimeCallback_cbS16IQ_LimeDevice_f36a06c4aed6a5ac(GoDeviceCallback *_swig_go_0, short *_swig_go_1, intgo _swig_go_2) {
  GoDeviceCallback *arg1 = (GoDeviceCallback *) 0 ;
  int16_t *arg2 = (int16_t *) 0 ;
  int arg3 ;
  
  arg1 = *(GoDeviceCallback **)&_swig_go_0; 
  arg2 = *(int16_t **)&_swig_go_1; 
  arg3 = (int)_swig_go_2; 
  
  (arg1)->cbS16IQ(arg2,arg3);
  
}


void _wrap_LimeCallback_cbS8IQ_LimeDevice_f36a06c4aed6a5ac(GoDeviceCallback *_swig_go_0, char *_swig_go_1, intgo _swig_go_2) {
  GoDeviceCallback *arg1 = (GoDeviceCallback *) 0 ;
  int8_t *arg2 = (int8_t *) 0 ;
  int arg3 ;
  
  arg1 = *(GoDeviceCallback **)&_swig_go_0; 
  arg2 = *(int8_t **)&_swig_go_1; 
  arg3 = (int)_swig_go_2; 
  
  (arg1)->cbS8IQ(arg2,arg3);
  
}


void _wrap_delete_LimeCallback_LimeDevice_f36a06c4aed6a5ac(GoDeviceCallback *_swig_go_0) {
  GoDeviceCallback *arg1 = (GoDeviceCallback *) 0 ;
  
  arg1 = *(GoDeviceCallback **)&_swig_go_0; 
  
  delete arg1;
  
}


GoDeviceCallback *_wrap_new_LimeCallback_LimeDevice_f36a06c4aed6a5ac() {
  GoDeviceCallback *result = 0 ;
  GoDeviceCallback *_swig_go_result;
  
  
  result = (GoDeviceCallback *)new GoDeviceCallback();
  *(GoDeviceCallback **)&_swig_go_result = (GoDeviceCallback *)result; 
  return _swig_go_result;
}


std::vector< unsigned int > *_wrap_new_Vector32u__SWIG_0_LimeDevice_f36a06c4aed6a5ac() {
  std::vector< uint32_t > *result = 0 ;
  std::vector< unsigned int > *_swig_go_result;
  
  
  result = (std::vector< uint32_t > *)new std::vector< uint32_t >();
  *(std::vector< uint32_t > **)&_swig_go_result = (std::vector< uint32_t > *)result; 
  return _swig_go_result;
}


std::vector< unsigned int > *_wrap_new_Vector32u__SWIG_1_LimeDevice_f36a06c4aed6a5ac(long long _swig_go_0) {
  std::vector< unsigned int >::size_type arg1 ;
  std::vector< uint32_t > *result = 0 ;
  std::vector< unsigned int > *_swig_go_result;
  
  arg1 = (size_t)_swig_go_0; 
  
  result = (std::vector< uint32_t > *)new std::vector< uint32_t >(arg1);
  *(std::vector< uint32_t > **)&_swig_go_result = (std::vector< uint32_t > *)result; 
  return _swig_go_result;
}


long long _wrap_Vector32u_size_LimeDevice_f36a06c4aed6a5ac(std::vector< unsigned int > *_swig_go_0) {
  std::vector< uint32_t > *arg1 = (std::vector< uint32_t > *) 0 ;
  std::vector< unsigned int >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< uint32_t > **)&_swig_go_0; 
  
  result = ((std::vector< uint32_t > const *)arg1)->size();
  _swig_go_result = result; 
  return _swig_go_result;
}


long long _wrap_Vector32u_capacity_LimeDevice_f36a06c4aed6a5ac(std::vector< unsigned int > *_swig_go_0) {
  std::vector< uint32_t > *arg1 = (std::vector< uint32_t > *) 0 ;
  std::vector< unsigned int >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< uint32_t > **)&_swig_go_0; 
  
  result = ((std::vector< uint32_t > const *)arg1)->capacity();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_Vector32u_reserve_LimeDevice_f36a06c4aed6a5ac(std::vector< unsigned int > *_swig_go_0, long long _swig_go_1) {
  std::vector< uint32_t > *arg1 = (std::vector< uint32_t > *) 0 ;
  std::vector< unsigned int >::size_type arg2 ;
  
  arg1 = *(std::vector< uint32_t > **)&_swig_go_0; 
  arg2 = (size_t)_swig_go_1; 
  
  (arg1)->reserve(arg2);
  
}


bool _wrap_Vector32u_isEmpty_LimeDevice_f36a06c4aed6a5ac(std::vector< unsigned int > *_swig_go_0) {
  std::vector< uint32_t > *arg1 = (std::vector< uint32_t > *) 0 ;
  bool result;
  bool _swig_go_result;
  
  arg1 = *(std::vector< uint32_t > **)&_swig_go_0; 
  
  result = (bool)((std::vector< uint32_t > const *)arg1)->empty();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_Vector32u_clear_LimeDevice_f36a06c4aed6a5ac(std::vector< unsigned int > *_swig_go_0) {
  std::vector< uint32_t > *arg1 = (std::vector< uint32_t > *) 0 ;
  
  arg1 = *(std::vector< uint32_t > **)&_swig_go_0; 
  
  (arg1)->clear();
  
}


void _wrap_Vector32u_add_LimeDevice_f36a06c4aed6a5ac(std::vector< unsigned int > *_swig_go_0, intgo _swig_go_1) {
  std::vector< uint32_t > *arg1 = (std::vector< uint32_t > *) 0 ;
  std::vector< unsigned int >::value_type *arg2 = 0 ;
  
  arg1 = *(std::vector< uint32_t > **)&_swig_go_0; 
  arg2 = (std::vector< unsigned int >::value_type *)&_swig_go_1; 
  
  (arg1)->push_back((std::vector< unsigned int >::value_type const &)*arg2);
  
}


intgo _wrap_Vector32u_get_LimeDevice_f36a06c4aed6a5ac(std::vector< unsigned int > *_swig_go_0, intgo _swig_go_1) {
  std::vector< uint32_t > *arg1 = (std::vector< uint32_t > *) 0 ;
  int arg2 ;
  std::vector< unsigned int >::value_type *result = 0 ;
  intgo _swig_go_result;
  
  arg1 = *(std::vector< uint32_t > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  
  try {
    result = (std::vector< unsigned int >::value_type *) &std_vector_Sl_uint32_t_Sg__get(arg1,arg2);
  }
  catch(std::out_of_range &_e) {
    _swig_gopanic((&_e)->what());
  }
  
  _swig_go_result = (unsigned int)*result; 
  return _swig_go_result;
}


void _wrap_Vector32u_set_LimeDevice_f36a06c4aed6a5ac(std::vector< unsigned int > *_swig_go_0, intgo _swig_go_1, intgo _swig_go_2) {
  std::vector< uint32_t > *arg1 = (std::vector< uint32_t > *) 0 ;
  int arg2 ;
  std::vector< unsigned int >::value_type *arg3 = 0 ;
  
  arg1 = *(std::vector< uint32_t > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  arg3 = (std::vector< unsigned int >::value_type *)&_swig_go_2; 
  
  try {
    std_vector_Sl_uint32_t_Sg__set(arg1,arg2,(unsigned int const &)*arg3);
  }
  catch(std::out_of_range &_e) {
    _swig_gopanic((&_e)->what());
  }
  
  
}


void _wrap_delete_Vector32u_LimeDevice_f36a06c4aed6a5ac(std::vector< unsigned int > *_swig_go_0) {
  std::vector< uint32_t > *arg1 = (std::vector< uint32_t > *) 0 ;
  
  arg1 = *(std::vector< uint32_t > **)&_swig_go_0; 
  
  delete arg1;
  
}


std::vector< float > *_wrap_new_Vector32f__SWIG_0_LimeDevice_f36a06c4aed6a5ac() {
  std::vector< float > *result = 0 ;
  std::vector< float > *_swig_go_result;
  
  
  result = (std::vector< float > *)new std::vector< float >();
  *(std::vector< float > **)&_swig_go_result = (std::vector< float > *)result; 
  return _swig_go_result;
}


std::vector< float > *_wrap_new_Vector32f__SWIG_1_LimeDevice_f36a06c4aed6a5ac(long long _swig_go_0) {
  std::vector< float >::size_type arg1 ;
  std::vector< float > *result = 0 ;
  std::vector< float > *_swig_go_result;
  
  arg1 = (size_t)_swig_go_0; 
  
  result = (std::vector< float > *)new std::vector< float >(arg1);
  *(std::vector< float > **)&_swig_go_result = (std::vector< float > *)result; 
  return _swig_go_result;
}


long long _wrap_Vector32f_size_LimeDevice_f36a06c4aed6a5ac(std::vector< float > *_swig_go_0) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  std::vector< float >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  
  result = ((std::vector< float > const *)arg1)->size();
  _swig_go_result = result; 
  return _swig_go_result;
}


long long _wrap_Vector32f_capacity_LimeDevice_f36a06c4aed6a5ac(std::vector< float > *_swig_go_0) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  std::vector< float >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  
  result = ((std::vector< float > const *)arg1)->capacity();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_Vector32f_reserve_LimeDevice_f36a06c4aed6a5ac(std::vector< float > *_swig_go_0, long long _swig_go_1) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  std::vector< float >::size_type arg2 ;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  arg2 = (size_t)_swig_go_1; 
  
  (arg1)->reserve(arg2);
  
}


bool _wrap_Vector32f_isEmpty_LimeDevice_f36a06c4aed6a5ac(std::vector< float > *_swig_go_0) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  bool result;
  bool _swig_go_result;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  
  result = (bool)((std::vector< float > const *)arg1)->empty();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_Vector32f_clear_LimeDevice_f36a06c4aed6a5ac(std::vector< float > *_swig_go_0) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  
  (arg1)->clear();
  
}


void _wrap_Vector32f_add_LimeDevice_f36a06c4aed6a5ac(std::vector< float > *_swig_go_0, float _swig_go_1) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  std::vector< float >::value_type *arg2 = 0 ;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  arg2 = (std::vector< float >::value_type *)&_swig_go_1; 
  
  (arg1)->push_back((std::vector< float >::value_type const &)*arg2);
  
}


float _wrap_Vector32f_get_LimeDevice_f36a06c4aed6a5ac(std::vector< float > *_swig_go_0, intgo _swig_go_1) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  int arg2 ;
  std::vector< float >::value_type *result = 0 ;
  float _swig_go_result;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  
  try {
    result = (std::vector< float >::value_type *) &std_vector_Sl_float_Sg__get(arg1,arg2);
  }
  catch(std::out_of_range &_e) {
    _swig_gopanic((&_e)->what());
  }
  
  _swig_go_result = (float)*result; 
  return _swig_go_result;
}


void _wrap_Vector32f_set_LimeDevice_f36a06c4aed6a5ac(std::vector< float > *_swig_go_0, intgo _swig_go_1, float _swig_go_2) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  int arg2 ;
  std::vector< float >::value_type *arg3 = 0 ;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  arg3 = (std::vector< float >::value_type *)&_swig_go_2; 
  
  try {
    std_vector_Sl_float_Sg__set(arg1,arg2,(float const &)*arg3);
  }
  catch(std::out_of_range &_e) {
    _swig_gopanic((&_e)->what());
  }
  
  
}


void _wrap_delete_Vector32f_LimeDevice_f36a06c4aed6a5ac(std::vector< float > *_swig_go_0) {
  std::vector< float > *arg1 = (std::vector< float > *) 0 ;
  
  arg1 = *(std::vector< float > **)&_swig_go_0; 
  
  delete arg1;
  
}


std::vector< short > *_wrap_new_Vector16i__SWIG_0_LimeDevice_f36a06c4aed6a5ac() {
  std::vector< int16_t > *result = 0 ;
  std::vector< short > *_swig_go_result;
  
  
  result = (std::vector< int16_t > *)new std::vector< int16_t >();
  *(std::vector< int16_t > **)&_swig_go_result = (std::vector< int16_t > *)result; 
  return _swig_go_result;
}


std::vector< short > *_wrap_new_Vector16i__SWIG_1_LimeDevice_f36a06c4aed6a5ac(long long _swig_go_0) {
  std::vector< short >::size_type arg1 ;
  std::vector< int16_t > *result = 0 ;
  std::vector< short > *_swig_go_result;
  
  arg1 = (size_t)_swig_go_0; 
  
  result = (std::vector< int16_t > *)new std::vector< int16_t >(arg1);
  *(std::vector< int16_t > **)&_swig_go_result = (std::vector< int16_t > *)result; 
  return _swig_go_result;
}


long long _wrap_Vector16i_size_LimeDevice_f36a06c4aed6a5ac(std::vector< short > *_swig_go_0) {
  std::vector< int16_t > *arg1 = (std::vector< int16_t > *) 0 ;
  std::vector< short >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< int16_t > **)&_swig_go_0; 
  
  result = ((std::vector< int16_t > const *)arg1)->size();
  _swig_go_result = result; 
  return _swig_go_result;
}


long long _wrap_Vector16i_capacity_LimeDevice_f36a06c4aed6a5ac(std::vector< short > *_swig_go_0) {
  std::vector< int16_t > *arg1 = (std::vector< int16_t > *) 0 ;
  std::vector< short >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< int16_t > **)&_swig_go_0; 
  
  result = ((std::vector< int16_t > const *)arg1)->capacity();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_Vector16i_reserve_LimeDevice_f36a06c4aed6a5ac(std::vector< short > *_swig_go_0, long long _swig_go_1) {
  std::vector< int16_t > *arg1 = (std::vector< int16_t > *) 0 ;
  std::vector< short >::size_type arg2 ;
  
  arg1 = *(std::vector< int16_t > **)&_swig_go_0; 
  arg2 = (size_t)_swig_go_1; 
  
  (arg1)->reserve(arg2);
  
}


bool _wrap_Vector16i_isEmpty_LimeDevice_f36a06c4aed6a5ac(std::vector< short > *_swig_go_0) {
  std::vector< int16_t > *arg1 = (std::vector< int16_t > *) 0 ;
  bool result;
  bool _swig_go_result;
  
  arg1 = *(std::vector< int16_t > **)&_swig_go_0; 
  
  result = (bool)((std::vector< int16_t > const *)arg1)->empty();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_Vector16i_clear_LimeDevice_f36a06c4aed6a5ac(std::vector< short > *_swig_go_0) {
  std::vector< int16_t > *arg1 = (std::vector< int16_t > *) 0 ;
  
  arg1 = *(std::vector< int16_t > **)&_swig_go_0; 
  
  (arg1)->clear();
  
}


void _wrap_Vector16i_add_LimeDevice_f36a06c4aed6a5ac(std::vector< short > *_swig_go_0, short _swig_go_1) {
  std::vector< int16_t > *arg1 = (std::vector< int16_t > *) 0 ;
  std::vector< short >::value_type *arg2 = 0 ;
  
  arg1 = *(std::vector< int16_t > **)&_swig_go_0; 
  arg2 = (std::vector< short >::value_type *)&_swig_go_1; 
  
  (arg1)->push_back((std::vector< short >::value_type const &)*arg2);
  
}


short _wrap_Vector16i_get_LimeDevice_f36a06c4aed6a5ac(std::vector< short > *_swig_go_0, intgo _swig_go_1) {
  std::vector< int16_t > *arg1 = (std::vector< int16_t > *) 0 ;
  int arg2 ;
  std::vector< short >::value_type *result = 0 ;
  short _swig_go_result;
  
  arg1 = *(std::vector< int16_t > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  
  try {
    result = (std::vector< short >::value_type *) &std_vector_Sl_int16_t_Sg__get(arg1,arg2);
  }
  catch(std::out_of_range &_e) {
    _swig_gopanic((&_e)->what());
  }
  
  _swig_go_result = (short)*result; 
  return _swig_go_result;
}


void _wrap_Vector16i_set_LimeDevice_f36a06c4aed6a5ac(std::vector< short > *_swig_go_0, intgo _swig_go_1, short _swig_go_2) {
  std::vector< int16_t > *arg1 = (std::vector< int16_t > *) 0 ;
  int arg2 ;
  std::vector< short >::value_type *arg3 = 0 ;
  
  arg1 = *(std::vector< int16_t > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  arg3 = (std::vector< short >::value_type *)&_swig_go_2; 
  
  try {
    std_vector_Sl_int16_t_Sg__set(arg1,arg2,(short const &)*arg3);
  }
  catch(std::out_of_range &_e) {
    _swig_gopanic((&_e)->what());
  }
  
  
}


void _wrap_delete_Vector16i_LimeDevice_f36a06c4aed6a5ac(std::vector< short > *_swig_go_0) {
  std::vector< int16_t > *arg1 = (std::vector< int16_t > *) 0 ;
  
  arg1 = *(std::vector< int16_t > **)&_swig_go_0; 
  
  delete arg1;
  
}


std::vector< signed char > *_wrap_new_Vector8i__SWIG_0_LimeDevice_f36a06c4aed6a5ac() {
  std::vector< int8_t > *result = 0 ;
  std::vector< signed char > *_swig_go_result;
  
  
  result = (std::vector< int8_t > *)new std::vector< int8_t >();
  *(std::vector< int8_t > **)&_swig_go_result = (std::vector< int8_t > *)result; 
  return _swig_go_result;
}


std::vector< signed char > *_wrap_new_Vector8i__SWIG_1_LimeDevice_f36a06c4aed6a5ac(long long _swig_go_0) {
  std::vector< signed char >::size_type arg1 ;
  std::vector< int8_t > *result = 0 ;
  std::vector< signed char > *_swig_go_result;
  
  arg1 = (size_t)_swig_go_0; 
  
  result = (std::vector< int8_t > *)new std::vector< int8_t >(arg1);
  *(std::vector< int8_t > **)&_swig_go_result = (std::vector< int8_t > *)result; 
  return _swig_go_result;
}


long long _wrap_Vector8i_size_LimeDevice_f36a06c4aed6a5ac(std::vector< signed char > *_swig_go_0) {
  std::vector< int8_t > *arg1 = (std::vector< int8_t > *) 0 ;
  std::vector< signed char >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< int8_t > **)&_swig_go_0; 
  
  result = ((std::vector< int8_t > const *)arg1)->size();
  _swig_go_result = result; 
  return _swig_go_result;
}


long long _wrap_Vector8i_capacity_LimeDevice_f36a06c4aed6a5ac(std::vector< signed char > *_swig_go_0) {
  std::vector< int8_t > *arg1 = (std::vector< int8_t > *) 0 ;
  std::vector< signed char >::size_type result;
  long long _swig_go_result;
  
  arg1 = *(std::vector< int8_t > **)&_swig_go_0; 
  
  result = ((std::vector< int8_t > const *)arg1)->capacity();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_Vector8i_reserve_LimeDevice_f36a06c4aed6a5ac(std::vector< signed char > *_swig_go_0, long long _swig_go_1) {
  std::vector< int8_t > *arg1 = (std::vector< int8_t > *) 0 ;
  std::vector< signed char >::size_type arg2 ;
  
  arg1 = *(std::vector< int8_t > **)&_swig_go_0; 
  arg2 = (size_t)_swig_go_1; 
  
  (arg1)->reserve(arg2);
  
}


bool _wrap_Vector8i_isEmpty_LimeDevice_f36a06c4aed6a5ac(std::vector< signed char > *_swig_go_0) {
  std::vector< int8_t > *arg1 = (std::vector< int8_t > *) 0 ;
  bool result;
  bool _swig_go_result;
  
  arg1 = *(std::vector< int8_t > **)&_swig_go_0; 
  
  result = (bool)((std::vector< int8_t > const *)arg1)->empty();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_Vector8i_clear_LimeDevice_f36a06c4aed6a5ac(std::vector< signed char > *_swig_go_0) {
  std::vector< int8_t > *arg1 = (std::vector< int8_t > *) 0 ;
  
  arg1 = *(std::vector< int8_t > **)&_swig_go_0; 
  
  (arg1)->clear();
  
}


void _wrap_Vector8i_add_LimeDevice_f36a06c4aed6a5ac(std::vector< signed char > *_swig_go_0, char _swig_go_1) {
  std::vector< int8_t > *arg1 = (std::vector< int8_t > *) 0 ;
  std::vector< signed char >::value_type *arg2 = 0 ;
  
  arg1 = *(std::vector< int8_t > **)&_swig_go_0; 
  arg2 = (std::vector< signed char >::value_type *)&_swig_go_1; 
  
  (arg1)->push_back((std::vector< signed char >::value_type const &)*arg2);
  
}


char _wrap_Vector8i_get_LimeDevice_f36a06c4aed6a5ac(std::vector< signed char > *_swig_go_0, intgo _swig_go_1) {
  std::vector< int8_t > *arg1 = (std::vector< int8_t > *) 0 ;
  int arg2 ;
  std::vector< signed char >::value_type *result = 0 ;
  char _swig_go_result;
  
  arg1 = *(std::vector< int8_t > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  
  try {
    result = (std::vector< signed char >::value_type *) &std_vector_Sl_int8_t_Sg__get(arg1,arg2);
  }
  catch(std::out_of_range &_e) {
    _swig_gopanic((&_e)->what());
  }
  
  _swig_go_result = (signed char)*result; 
  return _swig_go_result;
}


void _wrap_Vector8i_set_LimeDevice_f36a06c4aed6a5ac(std::vector< signed char > *_swig_go_0, intgo _swig_go_1, char _swig_go_2) {
  std::vector< int8_t > *arg1 = (std::vector< int8_t > *) 0 ;
  int arg2 ;
  std::vector< signed char >::value_type *arg3 = 0 ;
  
  arg1 = *(std::vector< int8_t > **)&_swig_go_0; 
  arg2 = (int)_swig_go_1; 
  arg3 = (std::vector< signed char >::value_type *)&_swig_go_2; 
  
  try {
    std_vector_Sl_int8_t_Sg__set(arg1,arg2,(signed char const &)*arg3);
  }
  catch(std::out_of_range &_e) {
    _swig_gopanic((&_e)->what());
  }
  
  
}


void _wrap_delete_Vector8i_LimeDevice_f36a06c4aed6a5ac(std::vector< signed char > *_swig_go_0) {
  std::vector< int8_t > *arg1 = (std::vector< int8_t > *) 0 ;
  
  arg1 = *(std::vector< int8_t > **)&_swig_go_0; 
  
  delete arg1;
  
}


LimeDevice *_wrap_new_LimeDevice_LimeDevice_f36a06c4aed6a5ac() {
  LimeDevice *result = 0 ;
  LimeDevice *_swig_go_result;
  
  
  result = (LimeDevice *)new LimeDevice();
  *(LimeDevice **)&_swig_go_result = (LimeDevice *)result; 
  return _swig_go_result;
}


void _wrap_delete_LimeDevice_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  
  delete arg1;
  
}


_gostring_ _wrap_LimeDevice_GetName_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  std::string *result = 0 ;
  _gostring_ _swig_go_result;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  
  result = (std::string *) &(arg1)->GetName();
  _swig_go_result = Swig_AllocateString((*result).data(), (*result).length()); 
  return _swig_go_result;
}


intgo _wrap_LimeDevice_SetSampleRate_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0, intgo _swig_go_1) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  uint32_t arg2 ;
  uint32_t result;
  intgo _swig_go_result;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  arg2 = (uint32_t)_swig_go_1; 
  
  result = (uint32_t)(arg1)->SetSampleRate(arg2);
  _swig_go_result = result; 
  return _swig_go_result;
}


intgo _wrap_LimeDevice_SetCenterFrequency_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0, intgo _swig_go_1) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  uint32_t arg2 ;
  uint32_t result;
  intgo _swig_go_result;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  arg2 = (uint32_t)_swig_go_1; 
  
  result = (uint32_t)(arg1)->SetCenterFrequency(arg2);
  _swig_go_result = result; 
  return _swig_go_result;
}


intgo _wrap_LimeDevice_GetCenterFrequency_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  uint32_t result;
  intgo _swig_go_result;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  
  result = (uint32_t)(arg1)->GetCenterFrequency();
  _swig_go_result = result; 
  return _swig_go_result;
}


intgo _wrap_LimeDevice_GetSampleRate_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  uint32_t result;
  intgo _swig_go_result;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  
  result = (uint32_t)(arg1)->GetSampleRate();
  _swig_go_result = result; 
  return _swig_go_result;
}


void _wrap_LimeDevice_SetLNAGain_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0, char _swig_go_1) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  uint8_t arg2 ;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  arg2 = (uint8_t)_swig_go_1; 
  
  (arg1)->SetLNAGain(arg2);
  
}


void _wrap_LimeDevice_SetTIAGain_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0, char _swig_go_1) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  uint8_t arg2 ;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  arg2 = (uint8_t)_swig_go_1; 
  
  (arg1)->SetTIAGain(arg2);
  
}


void _wrap_LimeDevice_SetPGAGain_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0, char _swig_go_1) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  uint8_t arg2 ;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  arg2 = (uint8_t)_swig_go_1; 
  
  (arg1)->SetPGAGain(arg2);
  
}


void _wrap_LimeDevice_SetAntenna_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0, _gostring_ _swig_go_1) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  std::string arg2 ;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  (&arg2)->assign(_swig_go_1.p, _swig_go_1.n); 
  
  (arg1)->SetAntenna(arg2);
  
}


void _wrap_LimeDevice_SetSamplesAvailableCallback_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0, GoDeviceCallback *_swig_go_1) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  GoDeviceCallback *arg2 = (GoDeviceCallback *) 0 ;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  arg2 = *(GoDeviceCallback **)&_swig_go_1; 
  
  (arg1)->SetSamplesAvailableCallback(arg2);
  
}


void _wrap_LimeDevice_Start_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  
  (arg1)->Start();
  
}


void _wrap_LimeDevice_Stop_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  
  (arg1)->Stop();
  
}


void _wrap_LimeDevice_GetSamples_LimeDevice_f36a06c4aed6a5ac(LimeDevice *_swig_go_0, short _swig_go_1) {
  LimeDevice *arg1 = (LimeDevice *) 0 ;
  uint16_t arg2 ;
  
  arg1 = *(LimeDevice **)&_swig_go_0; 
  arg2 = (uint16_t)_swig_go_1; 
  
  (arg1)->GetSamples(arg2);
  
}


#ifdef __cplusplus
}
#endif

